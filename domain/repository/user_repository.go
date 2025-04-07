package repository

import (
	"context"
	"record-project/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	FindByOpenID(ctx context.Context, openID string) (*entity.User, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint64) error

	// SearchUsers 搜索用户
	SearchUsers(ctx context.Context, keyword string, page, size int, excludeUserID uint64) ([]*entity.User, int64, error)
	
	// SearchUsersExcludeFriends 搜索用户（排除好友）
	SearchUsersExcludeFriends(ctx context.Context, keyword string, page, size int, excludeUserID uint64, friendIDs []uint64) ([]*entity.User, int64, error)
}
