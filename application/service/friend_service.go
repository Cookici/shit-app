package service

import (
	"context"
	"errors"
	"record-project/domain/entity"
	"record-project/domain/repository"
)

// FriendService 好友服务接口
type FriendService interface {
	// GetFriendsByUserID 获取用户的好友列表
	GetFriendsByUserID(ctx context.Context, userID uint64) ([]*entity.Friend, error)

	// GetFriendIDs 获取用户的好友ID列表
	GetFriendIDs(ctx context.Context, userID uint64) ([]uint64, error)

	// AddFriend 添加好友
	AddFriend(ctx context.Context, userID, friendID uint64) error

	// UpdateFriendStatus 更新好友状态
	UpdateFriendStatus(ctx context.Context, id uint64, status int8) error

	// DeleteFriend 删除好友
	DeleteFriend(ctx context.Context, friendId uint64, userID uint64) error

	// GetFriendRelation 获取好友关系
	GetFriendRelation(ctx context.Context, userID, friendID uint64) (*entity.Friend, error)

	// GetFriendsByUserIDWithPagination 分页获取用户的好友列表
	GetFriendsByUserIDWithPagination(ctx context.Context, userID uint64, page, size int, keyword string) ([]*entity.Friend, int64, error)

	// GetFriendRequests 获取发送给用户的好友申请
	GetFriendRequests(ctx context.Context, userID uint64, page, size int) ([]*entity.Friend, int64, error)

	// GetFriendRelationByID 通过ID获取好友关系
	GetFriendRelationByID(ctx context.Context, relationID uint64) (*entity.Friend, error)

	// UpdateFriendStatusWithVerification 更新好友状态（带验证）
	UpdateFriendStatusWithVerification(ctx context.Context, relationID, userID, friendID uint64, status int8) error
}

// friendService 好友服务实现
type friendService struct {
	friendRepo repository.FriendRepository
}

// NewFriendService 创建好友服务
func NewFriendService(friendRepo repository.FriendRepository) FriendService {
	return &friendService{
		friendRepo: friendRepo,
	}
}

// GetFriendsByUserID 获取用户的好友列表
func (s *friendService) GetFriendsByUserID(ctx context.Context, userID uint64) ([]*entity.Friend, error) {
	return s.friendRepo.FindByUserID(ctx, userID)
}

// GetFriendIDs 获取用户的好友ID列表
func (s *friendService) GetFriendIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	// 获取已确认状态的好友ID列表
	friendIDs, err := s.friendRepo.FindFriendIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return friendIDs, nil
}

// AddFriend 添加好友
func (s *friendService) AddFriend(ctx context.Context, userID, friendID uint64) error {
	friend := &entity.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   0, // 待确认状态
	}
	return s.friendRepo.Save(ctx, friend)
}

// GetFriendsByUserIDWithPagination 分页获取用户的好友列表
func (s *friendService) GetFriendsByUserIDWithPagination(ctx context.Context, userID uint64, page, size int, keyword string) ([]*entity.Friend, int64, error) {
	friends, total, err := s.friendRepo.FindByUserIDWithPagination(ctx, userID, page, size, keyword)
	if err != nil {
		return nil, 0, err
	}

	return friends, total, nil
}

// GetFriendRequests 获取发送给用户的好友申请
func (s *friendService) GetFriendRequests(ctx context.Context, userID uint64, page, size int) ([]*entity.Friend, int64, error) {
	friends, total, err := s.friendRepo.FindFriendRequestsByUserID(ctx, userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	return friends, total, nil
}

// UpdateFriendStatus 更新好友状态
func (s *friendService) UpdateFriendStatus(ctx context.Context, id uint64, status int8) error {
	// 查找好友关系
	var friend *entity.Friend
	var err error

	// 先尝试通过ID查找
	friends, err := s.friendRepo.FindByUserID(ctx, id)
	if err != nil {
		return err
	}

	if len(friends) > 0 {
		friend = friends[0]
	}

	if friend == nil {
		return errors.New("好友关系不存在")
	}

	friend.Status = status
	return s.friendRepo.Update(ctx, friend)
}

// DeleteFriend 删除好友
func (s *friendService) DeleteFriend(ctx context.Context, id uint64, userID uint64) error {
	return s.friendRepo.Delete(ctx, id, userID)
}

// GetFriendRelation 获取好友关系
func (s *friendService) GetFriendRelation(ctx context.Context, userID, friendID uint64) (*entity.Friend, error) {
	return s.friendRepo.FindByUserIDAndFriendID(ctx, userID, friendID)
}

// GetFriendRelationByID 通过ID获取好友关系
func (s *friendService) GetFriendRelationByID(ctx context.Context, relationID uint64) (*entity.Friend, error) {
	return s.friendRepo.FindByID(ctx, relationID)
}

// UpdateFriendStatusWithVerification 更新好友状态（带验证）
func (s *friendService) UpdateFriendStatusWithVerification(ctx context.Context, relationID, userID, friendID uint64, status int8) error {
	// 获取好友关系
	relation, err := s.friendRepo.FindByID(ctx, relationID)
	if err != nil {
		return err
	}

	if relation == nil {
		return errors.New("好友关系不存在")
	}
	
	// 验证用户权限
	if relation.UserID != userID && relation.FriendID != userID {
		return errors.New("无权更新此好友关系")
	}

	// 验证好友ID
	if (relation.UserID == userID && relation.FriendID != friendID) ||
		(relation.FriendID == userID && relation.UserID != friendID) {
		return errors.New("好友ID不匹配")
	}

	// 更新状态
	relation.Status = status
	return s.friendRepo.Update(ctx, relation)
}
