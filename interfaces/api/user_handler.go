package api

import (
	"net/http"
	"record-project/application/service"
	"record-project/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户API处理器
type UserHandler struct {
	userService   service.UserService
	authService   service.AuthService
	friendService service.FriendService // 添加好友服务
}

// NewUserHandler 创建用户API处理器
func NewUserHandler(userService service.UserService, authService service.AuthService, friendService service.FriendService) *UserHandler {
	return &UserHandler{
		userService:   userService,
		authService:   authService,
		friendService: friendService,
	}
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := h.userService.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByOpenID 根据OpenID获取用户
func (h *UserHandler) GetUserByOpenID(c *gin.Context) {
	openID := c.Param("openid")
	if openID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OpenID不能为空"})
		return
	}

	user, err := h.userService.GetUserByOpenID(c, openID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := h.userService.UpdateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := h.userService.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// SearchUsers 搜索用户（结果不包含已是好友的用户）
func (h *UserHandler) SearchUsers(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	userID, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取查询参数
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// 获取当前用户的好友ID列表
	friendIDs, err := h.friendService.GetFriendIDs(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友列表失败: " + err.Error()})
		return
	}

	// 调用服务层搜索用户，排除好友
	users, total, err := h.userService.SearchUsersExcludeFriends(c, keyword, page, pageSize, userID, friendIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索用户失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
