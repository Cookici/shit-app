package service

import (
	"context"
	"record-project/domain/entity"
	"record-project/domain/repository"
)

// UserService 用户服务接口
type UserService interface {
	GetUserByID(ctx context.Context, id uint64) (*entity.User, error)
	GetUserByOpenID(ctx context.Context, openID string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uint64) error
	GetUsersByIDs(ctx context.Context, ids []uint64) ([]*entity.User, error)

	// SearchUsers 搜索用户
	SearchUsers(ctx context.Context, keyword string, page, size int, currentUserID uint64) ([]*entity.User, int, error)

	// SearchUsersExcludeFriends 搜索用户（排除好友）
	SearchUsersExcludeFriends(ctx context.Context, keyword string, page, size int, currentUserID uint64, friendIDs []uint64) ([]*entity.User, int64, error)
}

// GetUsersByIDs 在userService实现中添加
func (s *userService) GetUsersByIDs(ctx context.Context, ids []uint64) ([]*entity.User, error) {
	return s.userRepo.FindByIDs(ctx, ids)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(ctx context.Context, id uint64) (*entity.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetUserByOpenID 根据OpenID获取用户
func (s *userService) GetUserByOpenID(ctx context.Context, openID string) (*entity.User, error) {
	return s.userRepo.FindByOpenID(ctx, openID)
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	return s.userRepo.Save(ctx, user)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint64) error {
	return s.userRepo.Delete(ctx, id)
}

// SearchUsers 搜索用户
func (s *userService) SearchUsers(ctx context.Context, keyword string, page, size int, currentUserID uint64) ([]*entity.User, int, error) {
	users, total, err := s.userRepo.SearchUsers(ctx, keyword, page, size, currentUserID)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int((total + int64(size) - 1) / int64(size)) // 计算总页数
	if totalPages < 1 {
		totalPages = 1
	}

	return users, totalPages, nil
}

// SearchUsersExcludeFriends 搜索用户（排除好友）
func (s *userService) SearchUsersExcludeFriends(ctx context.Context, keyword string, page, size int, currentUserID uint64, friendIDs []uint64) ([]*entity.User, int64, error) {
	users, total, err := s.userRepo.SearchUsersExcludeFriends(ctx, keyword, page, size, currentUserID, friendIDs)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
