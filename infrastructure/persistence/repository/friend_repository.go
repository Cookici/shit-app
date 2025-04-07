package repository

import (
	"context"
	"errors"
	"fmt"
	"record-project/domain/entity"
	"record-project/domain/repository"
	"record-project/infrastructure/persistence/model"
	"time"

	"gorm.io/gorm"
)

// friendRepository 好友关系仓储实现
type friendRepository struct {
	db *gorm.DB
}

// NewFriendRepository 创建好友关系仓储
func NewFriendRepository(db *gorm.DB) repository.FriendRepository {
	return &friendRepository{db: db}
}

// FindByUserID 查找用户的好友列表
func (r *friendRepository) FindByUserID(ctx context.Context, userID uint64) ([]*entity.Friend, error) {
	var friendModels []model.Friend
	// 同时查询user_id和friend_id字段
	if err := r.db.WithContext(ctx).Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, 1).Find(&friendModels).Error; err != nil {
		return nil, err
	}

	friends := make([]*entity.Friend, len(friendModels))
	for i, friendModel := range friendModels {
		friends[i] = friendModel.ToEntity()

		// 确定实际的好友ID
		var actualFriendID uint64
		if friendModel.UserID == userID {
			actualFriendID = friendModel.FriendID
		} else {
			actualFriendID = friendModel.UserID
		}

		// 查询好友的用户信息
		var userModel model.User
		if err := r.db.WithContext(ctx).First(&userModel, actualFriendID).Error; err == nil {
			friends[i].FriendUser = userModel.ToEntity()
			// 确保FriendID字段始终是好友的ID
			friends[i].FriendID = actualFriendID
		}
	}

	return friends, nil
}

// FindByUserIDWithPagination 分页查询用户的好友列表
func (r *friendRepository) FindByUserIDWithPagination(ctx context.Context, userID uint64, page, size int, keyword string) ([]*entity.Friend, int64, error) {
	var friendModels []model.Friend
	var total int64

	// 构建查询 - 同时查询user_id和friend_id字段
	query := r.db.WithContext(ctx).Model(&model.Friend{}).Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, 1)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 如果有关键词，需要联表查询用户表
	if keyword != "" {
		// 使用子查询获取符合条件的好友ID
		subQuery := r.db.Table("users").
			Select("id").
			Where("nickname LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

		// 修改查询条件，考虑双向关系
		query = query.Where("(user_id IN (?) AND friend_id = ?) OR (friend_id IN (?) AND user_id = ?)",
			subQuery, userID, subQuery, userID)

		// 重新计算总数
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&friendModels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为实体
	friends := make([]*entity.Friend, len(friendModels))
	for i, friendModel := range friendModels {
		friends[i] = friendModel.ToEntity()

		// 确定实际的好友ID
		var actualFriendID uint64
		if friendModel.UserID == userID {
			actualFriendID = friendModel.FriendID
		} else {
			actualFriendID = friendModel.UserID
		}

		// 查询好友的用户信息
		var userModel model.User
		if err := r.db.WithContext(ctx).First(&userModel, actualFriendID).Error; err == nil {
			friends[i].FriendUser = userModel.ToEntity()
			// 确保FriendID字段始终是好友的ID
			friends[i].FriendID = actualFriendID
		}
	}

	return friends, total, nil
}

// FindFriendIDs 获取用户的好友ID列表(已确认的好友)
func (r *friendRepository) FindFriendIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	var friendIDs []uint64

	// 查询已确认状态(status=1)的好友ID - 考虑双向关系
	if err := r.db.WithContext(ctx).Model(&model.Friend{}).
		Where("user_id = ? AND status = ?", userID, 1).
		Pluck("friend_id", &friendIDs).Error; err != nil {
		return nil, err
	}

	// 查询反向关系
	var reverseIDs []uint64
	if err := r.db.WithContext(ctx).Model(&model.Friend{}).
		Where("friend_id = ? AND status = ?", userID, 1).
		Pluck("user_id", &reverseIDs).Error; err != nil {
		return nil, err
	}

	// 合并两个列表
	friendIDs = append(friendIDs, reverseIDs...)

	return friendIDs, nil
}

// Save 保存好友关系
func (r *friendRepository) Save(ctx context.Context, friend *entity.Friend) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先检查是否已存在好友关系
		var existingFriend model.Friend
		err := tx.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			friend.UserID, friend.FriendID, friend.FriendID, friend.UserID).First(&existingFriend).Error

		if err == nil {
			// 已存在关系，根据状态处理
			if existingFriend.Status == 1 {
				return errors.New("已经是好友关系")
			} else if existingFriend.Status == 0 {
				// 如果是待确认状态，可以更新为已确认
				if existingFriend.UserID == friend.FriendID && existingFriend.FriendID == friend.UserID {
					// 对方发起的申请，我方同意
					existingFriend.Status = 1
					return tx.Model(&existingFriend).Update("status", 1).Error
				} else {
					return errors.New("已发送过好友申请，等待对方确认")
				}
			} else if existingFriend.Status == 2 {
				// 如果是已拒绝状态，可以重新申请
				existingFriend.Status = 0
				return tx.Model(&existingFriend).Update("status", 0).Error
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 不存在关系，创建新的
		var friendModel model.Friend
		friendModel.FromEntity(friend)
		if err := tx.Create(&friendModel).Error; err != nil {
			return err
		}
		friend.ID = friendModel.ID
		return nil
	})
}

// Update 更新好友状态
func (r *friendRepository) Update(ctx context.Context, friend *entity.Friend) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新状态
		result := tx.Model(&model.Friend{}).Where("id = ?", friend.ID).Updates(map[string]interface{}{
			"status":     friend.Status,
			"updated_at": time.Now(),
		})
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

// Delete 删除好友关系
func (r *friendRepository) Delete(ctx context.Context, id uint64, userID uint64) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除正向关系
		result := tx.Where("id = ? and ( user_id = ? or friend_id = ? )", id, userID, userID).Delete(&model.Friend{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("好友关系不存在")
		}

		return nil
	})
}

// FindByUserIDAndFriendID 查找特定的好友关系
func (r *friendRepository) FindByUserIDAndFriendID(ctx context.Context, userID, friendID uint64) (*entity.Friend, error) {
	var friendModel model.Friend
	// 查询正向关系
	err := r.db.WithContext(ctx).Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friendModel).Error
	if err == nil {
		return friendModel.ToEntity(), nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 查询反向关系
	err = r.db.WithContext(ctx).Where("user_id = ? AND friend_id = ?", friendID, userID).First(&friendModel).Error
	if err == nil {
		return friendModel.ToEntity(), nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}

// FindFriendRequestsByUserID 查询发送给用户的好友申请
func (r *friendRepository) FindFriendRequestsByUserID(ctx context.Context, userID uint64, page, size int) ([]*entity.Friend, int64, error) {
	var friendModels []model.Friend
	var total int64

	// 构建查询 - 查询发送给用户的待确认申请
	query := r.db.WithContext(ctx).Model(&model.Friend{}).
		Where("friend_id = ? AND status = ?", userID, 0) // status=0 表示待确认状态

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&friendModels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为实体
	friends := make([]*entity.Friend, len(friendModels))
	for i, friendModel := range friendModels {
		friends[i] = friendModel.ToEntity()

		// 查询申请者的用户信息
		var userModel model.User
		if err := r.db.WithContext(ctx).First(&userModel, friendModel.UserID).Error; err == nil {
			friends[i].FriendUser = userModel.ToEntity()
		}
	}

	return friends, total, nil
}

// FindByID 通过ID查找好友关系
func (r *friendRepository) FindByID(ctx context.Context, id uint64) (*entity.Friend, error) {
	var friendModel model.Friend
	err := r.db.WithContext(ctx).First(&friendModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return friendModel.ToEntity(), nil
}
