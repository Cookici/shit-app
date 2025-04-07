package api

import (
	"net/http"
	"record-project/application/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FriendHandler 好友API处理器
type FriendHandler struct {
	friendService service.FriendService
	authService   service.AuthService
	userService   service.UserService
}

// NewFriendHandler 创建好友API处理器
func NewFriendHandler(friendService service.FriendService, authService service.AuthService, userService service.UserService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
		authService:   authService,
		userService:   userService,
	}
}

// GetFriends 获取好友列表
func (h *FriendHandler) GetFriends(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	userID, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 10
	keyword := c.Query("keyword")

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

	// 获取好友列表
	friends, total, err := h.friendService.GetFriendsByUserIDWithPagination(c, userID, page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"friends":   friends,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AddFriend 添加好友
func (h *FriendHandler) AddFriend(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	userID, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 解析请求体
	var request struct {
		FriendID uint64 `json:"friend_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 不能添加自己为好友
	if request.FriendID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能添加自己为好友"})
		return
	}

	// 检查好友是否存在
	friend, err := h.userService.GetUserByID(c, request.FriendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	if friend == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 添加好友
	if err := h.friendService.AddFriend(c, userID, request.FriendID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "添加好友失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友请求已发送"})
}

// UpdateFriendStatus 更新好友状态
func (h *FriendHandler) UpdateFriendStatus(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	_, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 解析请求体
	var request struct {
		UserID   uint64 `json:"user_id" binding:"required"`
		FriendID uint64 `json:"friend_id" binding:"required"`
		Status   int8   `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 验证状态值
	if request.Status < 0 || request.Status > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态值"})
		return
	}

	// 获取好友关系ID
	relationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的好友关系ID"})
		return
	}

	// 更新好友状态
	if err := h.friendService.UpdateFriendStatusWithVerification(c, relationID, request.UserID, request.FriendID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新好友状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友状态已更新"})
}

// DeleteFriend 删除好友
func (h *FriendHandler) DeleteFriend(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	userID, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 解析请求体
	var request struct {
		ID uint64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 删除好友
	if err := h.friendService.DeleteFriend(c, request.ID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除好友失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友已删除"})
}

// GetFriendRequests 获取好友申请
func (h *FriendHandler) GetFriendRequests(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	userID, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
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

	// 调用服务层获取好友申请
	requests, total, err := h.friendService.GetFriendRequests(c, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友申请失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests":  requests,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
