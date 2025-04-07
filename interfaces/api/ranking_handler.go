package api

import (
	"net/http"
	"record-project/application/service"
	"record-project/domain/entity"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RankingHandler 排行榜API处理器
type RankingHandler struct {
	recordService service.RecordService
	authService   service.AuthService
	userService   service.UserService
	friendService service.FriendService
}

// NewRankingHandler 创建排行榜API处理器
func NewRankingHandler(recordService service.RecordService, authService service.AuthService, userService service.UserService, friendService service.FriendService) *RankingHandler {
	return &RankingHandler{
		recordService: recordService,
		authService:   authService,
		userService:   userService,
		friendService: friendService,
	}
}

// GetRanking 获取全局排行榜
func (h *RankingHandler) GetRanking(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	_, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取查询参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// 默认查询参数
	if startDateStr == "" {
		// 默认为当前月份的第一天
		now := time.Now()
		startDateStr = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}

	if endDateStr == "" {
		// 默认为今天
		endDateStr = time.Now().Format("2006-01-02")
	}

	// 解析日期并设置为东八区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	endDate, err := time.ParseInLocation("2006-01-02", endDateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	// 设置结束日期为当天的23:59:59（东八区）
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)

	// 获取全局排行榜数据（前10名）
	rankingItems, err := h.recordService.GetGlobalRanking(c, startDate, endDate, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜数据失败"})
		return
	}

	// 如果没有排行数据，直接返回空列表
	if len(rankingItems) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"rankings": []interface{}{},
		})
		return
	}

	// 收集所有需要查询的用户ID
	var userIDs []uint64
	for _, item := range rankingItems {
		userIDs = append(userIDs, item.UserID)
	}

	// 查询用户信息
	users, err := h.userService.GetUsersByIDs(c, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 创建用户ID到用户信息的映射
	userMap := make(map[uint64]*entity.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 组装完整的排行榜数据
	for _, item := range rankingItems {
		if user, exists := userMap[item.UserID]; exists {
			item.Nickname = user.Nickname
			item.AvatarURL = user.AvatarURL
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings": rankingItems,
	})
}

// GetFriendRanking 获取好友排行榜
func (h *RankingHandler) GetFriendRanking(c *gin.Context) {
	// 从请求中获取token，解析用户ID
	userID, err := h.authService.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取查询参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 默认查询参数
	if startDateStr == "" {
		// 默认为当前月份的第一天
		now := time.Now()
		startDateStr = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}

	if endDateStr == "" {
		// 默认为今天
		endDateStr = time.Now().Format("2006-01-02")
	}

	// 解析日期并设置为东八区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	endDate, err := time.ParseInLocation("2006-01-02", endDateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	// 设置结束日期为当天的23:59:59（东八区）
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)

	// 获取用户的好友列表
	friends, err := h.friendService.GetFriendsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友列表失败"})
		return
	}

	// 收集所有好友ID，包括自己
	friendIDs := []uint64{userID} // 包含自己
	for _, friend := range friends {
		if friend.FriendID != userID {
			friendIDs = append(friendIDs, friend.FriendID)
		}
	}

	// 获取好友排行榜数据
	rankingItems, total, err := h.recordService.GetFriendRanking(c, friendIDs, startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜数据失败"})
		return
	}

	// 收集所有需要查询的用户ID
	var userIDs []uint64
	for _, item := range rankingItems {
		userIDs = append(userIDs, item.UserID)
	}

	// 查询用户信息
	users, err := h.userService.GetUsersByIDs(c, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 创建用户ID到用户信息的映射
	userMap := make(map[uint64]*entity.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 组装完整的排行榜数据
	for _, item := range rankingItems {
		if user, exists := userMap[item.UserID]; exists {
			item.Nickname = user.Nickname
			item.AvatarURL = user.AvatarURL
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings":  rankingItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
