package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"record-project/domain/entity"
	"record-project/infrastructure/auth"
	"record-project/infrastructure/wechat"
	"time"
)

// AuthService 认证服务接口
// 确保AuthService接口中有以下方法
type AuthService interface {
	WechatLogin(ctx context.Context, code string) (*entity.User, string, error)
	UpdateUserInfo(ctx context.Context, userID uint64, nickname, avatarURL string) error
	GetUserIDFromToken(ctx *gin.Context) (uint64, error)
}

func (s *authService) GetUserIDFromToken(ctx *gin.Context) (uint64, error) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		return 0, errors.New("未找到用户ID")
	}

	// 转换为uint64类型
	id, ok := userID.(uint64)
	if !ok {
		return 0, errors.New("用户ID类型错误")
	}

	return id, nil
}

// authService 认证服务实现
type authService struct {
	userService   UserService
	wechatService wechat.WechatService
}

// NewAuthService 创建认证服务
func NewAuthService(userService UserService, wechatService wechat.WechatService) AuthService {
	return &authService{
		userService:   userService,
		wechatService: wechatService,
	}
}

// WechatLogin 微信登录
// 修改WechatLogin方法，使用JWT生成token
func (s *authService) WechatLogin(ctx context.Context, code string) (*entity.User, string, error) {
	// 调用微信服务获取openid
	wxResp, err := s.wechatService.Code2Session(code)
	if err != nil {
		return nil, "", err
	}

	// 查找用户是否存在
	user, err := s.userService.GetUserByOpenID(ctx, wxResp.OpenID)
	if err != nil {
		return nil, "", err
	}

	// 用户不存在则创建新用户
	if user == nil {
		openIDSuffix := wxResp.OpenID
		if len(openIDSuffix) > 6 {
			openIDSuffix = openIDSuffix[len(openIDSuffix)-6:]
		}

		user = &entity.User{
			OpenID:    wxResp.OpenID,
			Nickname:  fmt.Sprintf("用户_%s", openIDSuffix),
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.userService.CreateUser(ctx, user); err != nil {
			return nil, "", err
		}
	}

	// 使用JWT生成token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// UpdateUserInfo 更新用户信息
func (s *authService) UpdateUserInfo(ctx context.Context, userID uint64, nickname, avatarURL string) error {
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	user.Nickname = nickname
	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	return s.userService.UpdateUser(ctx, user)
}
