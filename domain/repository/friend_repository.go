package repository

import (
	"context"
	"record-project/domain/entity"
)

// FriendRepository 好友关系仓储接口
type FriendRepository interface {
	// FindByUserID 查找用户的好友列表
	FindByUserID(ctx context.Context, userID uint64) ([]*entity.Friend, error)

	// FindFriendIDs 获取用户的好友ID列表(已确认的好友)
	FindFriendIDs(ctx context.Context, userID uint64) ([]uint64, error)

	// Save 保存好友关系
	Save(ctx context.Context, friend *entity.Friend) error

	// Update 更新好友关系
	Update(ctx context.Context, friend *entity.Friend) error

	// Delete 删除好友关系
	Delete(ctx context.Context, id uint64, userID uint64) error

	// FindByUserIDAndFriendID 查找特定的好友关系
	FindByUserIDAndFriendID(ctx context.Context, userID, friendID uint64) (*entity.Friend, error)

	// FindByUserIDWithPagination 分页查询用户的好友列表
	FindByUserIDWithPagination(ctx context.Context, userID uint64, page, size int, keyword string) ([]*entity.Friend, int64, error)

	// FindFriendRequestsByUserID 查询发送给用户的好友申请
	FindFriendRequestsByUserID(ctx context.Context, userID uint64, page, size int) ([]*entity.Friend, int64, error)

	// FindByID 通过ID查找好友关系
	FindByID(ctx context.Context, id uint64) (*entity.Friend, error)
}
