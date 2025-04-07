package repository

import (
	"context"
	"errors"
	"record-project/domain/entity"
	"record-project/domain/repository"
	"record-project/infrastructure/persistence/model"
	"time"

	"gorm.io/gorm"
)

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userModel.ToEntity(), nil
}

// FindByOpenID 根据OpenID查找用户
func (r *userRepository) FindByOpenID(ctx context.Context, openID string) (*entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).Where("open_id = ?", openID).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userModel.ToEntity(), nil
}

// FindByIDs 根据ID列表查找用户
func (r *userRepository) FindByIDs(ctx context.Context, ids []uint64) ([]*entity.User, error) {
	var userModels []model.User
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&userModels).Error; err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = userModel.ToEntity()
	}

	return users, nil
}

// Save 保存用户
func (r *userRepository) Save(ctx context.Context, user *entity.User) error {
	var userModel model.User
	userModel.FromEntity(user)
	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return err
	}
	user.ID = userModel.ID
	return nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	var userModel model.User
	userModel.FromEntity(user)

	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"open_id":    userModel.OpenID,
		"nickname":   userModel.Nickname,
		"avatar_url": userModel.AvatarURL,
		"status":     userModel.Status,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// SearchUsers 搜索用户
func (r *userRepository) SearchUsers(ctx context.Context, keyword string, page, size int, excludeUserID uint64) ([]*entity.User, int64, error) {
	var userModels []model.User
	var total int64

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id != ?", excludeUserID). // 排除当前用户
		Where("nickname LIKE ?", "%"+keyword+"%")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&userModels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为实体
	users := make([]*entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = userModel.ToEntity()
	}

	return users, total, nil
}

// SearchUsersExcludeFriends 搜索用户（排除好友）
func (r *userRepository) SearchUsersExcludeFriends(ctx context.Context, keyword string, page, size int, excludeUserID uint64, friendIDs []uint64) ([]*entity.User, int64, error) {
	var userModels []model.User
	var total int64

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id != ?", excludeUserID) // 排除当前用户

	// 如果有好友ID列表，排除这些ID
	if len(friendIDs) > 0 {
		query = query.Where("id NOT IN (?)", friendIDs)
	}

	// 添加关键词搜索条件
	query = query.Where("nickname LIKE ?", "%"+keyword+"%")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&userModels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为实体
	users := make([]*entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = userModel.ToEntity()
	}

	return users, total, nil
}
