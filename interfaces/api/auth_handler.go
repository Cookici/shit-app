package api

import (
	"net/http"
	"record-project/application/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证API处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证API处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// WechatLogin 微信登录
func (h *AuthHandler) WechatLogin(c *gin.Context) {
	var request struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	user, token, err := h.authService.WechatLogin(c, request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// UpdateUserInfo 更新用户信息
func (h *AuthHandler) UpdateUserInfo(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var request struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.UpdateUserInfo(c, userID, request.Nickname, request.AvatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}
