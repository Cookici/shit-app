package api

import (
	"net/http"
	"record-project/application/service"
	"record-project/domain/entity"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RecordHandler 记录API处理器
type RecordHandler struct {
	recordService   service.RecordService
	userService     service.UserService
	tagService      service.TagService
	poopTypeService service.PoopTypeService
}

// NewRecordHandler 创建记录API处理器
func NewRecordHandler(
	recordService service.RecordService,
	userService service.UserService,
	tagService service.TagService,
	poopTypeService service.PoopTypeService,
) *RecordHandler {
	return &RecordHandler{
		recordService:   recordService,
		userService:     userService,
		tagService:      tagService,
		poopTypeService: poopTypeService,
	}
}

// GetRecord 获取记录
func (h *RecordHandler) GetRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	record, err := h.recordService.GetRecordByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	// 获取关联的标签
	tags, err := h.recordService.GetRecordTags(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取关联的用户
	user, err := h.userService.GetUserByID(c, record.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取关联的屎的类型
	poopType, err := h.poopTypeService.GetPoopTypeByID(c, record.PoopTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"record":    record,
		"tags":      tags,
		"user":      user,
		"poop_type": poopType,
	}

	c.JSON(http.StatusOK, response)
}

// GetRecordsByUserID 根据用户ID获取记录列表
func (h *RecordHandler) GetRecordsByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	records, total, err := h.recordService.GetRecordsByUserID(c, userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 批量处理记录关联数据
	if len(records) > 0 {
		// 1. 收集所有记录ID和屎类型ID
		var recordIDs []uint64
		var poopTypeIDs []uint64
		recordIDMap := make(map[uint64]*entity.Record)

		for _, record := range records {
			recordIDs = append(recordIDs, record.ID)
			if record.PoopTypeID > 0 {
				poopTypeIDs = append(poopTypeIDs, record.PoopTypeID)
			}
			recordIDMap[record.ID] = record
		}

		// 2. 批量获取所有记录的标签
		recordTagsMap := make(map[uint64][]*entity.Tag)
		allTags, err := h.tagService.GetTagsByRecordIDs(c, recordIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for recordID, tags := range allTags {
			recordTagsMap[recordID] = tags
		}

		// 3. 批量获取所有屎类型
		poopTypeMap := make(map[uint64]*entity.PoopType)
		if len(poopTypeIDs) > 0 {
			poopTypes, err := h.poopTypeService.GetPoopTypesByIDs(c, poopTypeIDs)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, pt := range poopTypes {
				poopTypeMap[pt.ID] = pt
			}
		}

		// 4. 组装完整记录
		var completeRecords []map[string]interface{}
		for _, record := range records {
			completeRecord := map[string]interface{}{
				"record":    record,
				"tags":      recordTagsMap[record.ID],
				"poop_type": poopTypeMap[record.PoopTypeID],
			}
			completeRecords = append(completeRecords, completeRecord)
		}

		c.JSON(http.StatusOK, gin.H{
			"records": completeRecords,
			"total":   total,
			"page":    page,
			"size":    size,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"records": []map[string]interface{}{},
			"total":   total,
			"page":    page,
			"size":    size,
		})
	}
}

// GetRecordsByDateRange 根据日期范围获取记录
func (h *RecordHandler) GetRecordsByDateRange(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))

	startStr := c.Query("start")
	endStr := c.Query("end")

	// 定义中国时区
	cst := time.FixedZone("CST", 8*3600)

	// 解析开始日期
	startTime, err := time.ParseInLocation("2006-01-02", startStr, cst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	// 解析结束日期
	endTime, err := time.ParseInLocation("2006-01-02", endStr, cst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	// 设置结束日期为当天的23:59:59
	endTime = endTime.Add(24*time.Hour - time.Second)

	total, err := h.recordService.CountRecordsByDateRange(c, userID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if total == 0 {
		c.JSON(http.StatusOK, gin.H{
			"records": []map[string]interface{}{},
			"total":   0,
		})
		return
	}

	// 直接传递time.Time对象，而不是字符串
	records, err := h.recordService.GetRecordsByDateRange(c, userID, startTime, endTime, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 1. 收集所有记录ID和屎类型ID
	var recordIDs []uint64
	var poopTypeIDs []uint64
	recordIDMap := make(map[uint64]*entity.Record)

	for _, record := range records {
		recordIDs = append(recordIDs, record.ID)
		if record.PoopTypeID > 0 {
			poopTypeIDs = append(poopTypeIDs, record.PoopTypeID)
		}
		recordIDMap[record.ID] = record
	}

	// 2. 批量获取所有记录的标签
	recordTagsMap := make(map[uint64][]*entity.Tag)
	allTags, err := h.tagService.GetTagsByRecordIDs(c, recordIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for recordID, tags := range allTags {
		recordTagsMap[recordID] = tags
	}

	// 3. 批量获取所有屎类型
	poopTypeMap := make(map[uint64]*entity.PoopType)
	if len(poopTypeIDs) > 0 {
		poopTypes, err := h.poopTypeService.GetPoopTypesByIDs(c, poopTypeIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, pt := range poopTypes {
			poopTypeMap[pt.ID] = pt
		}
	}

	// 4. 组装完整记录
	var completeRecords []map[string]interface{}
	for _, record := range records {
		completeRecord := map[string]interface{}{
			"record":    record,
			"tags":      recordTagsMap[record.ID],
			"poop_type": poopTypeMap[record.PoopTypeID],
		}
		completeRecords = append(completeRecords, completeRecord)
	}

	c.JSON(http.StatusOK, gin.H{
		"records": completeRecords,
		"total":   total,
	})
}

// CreateRecord 创建记录
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var request struct {
		Record *entity.Record `json:"record"`
		TagIDs []uint64       `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用事务创建记录并关联标签
	if err := h.recordService.CreateRecordWithTags(c, request.Record, request.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request.Record)
}

// UpdateRecord 更新记录
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var request struct {
		Record *entity.Record `json:"record"`
		TagIDs []uint64       `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Record.ID = id
	// 使用事务更新记录并关联标签
	if err := h.recordService.UpdateRecordWithTags(c, request.Record, request.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request.Record)
}

// DeleteRecord 删除记录
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	if err := h.recordService.DeleteRecord(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "记录删除成功"})
}

// GetUsersDailyRecordStats 批量获取指定用户当天的拉屎记录统计
func (h *RecordHandler) GetUsersDailyRecordStats(c *gin.Context) {
	// 获取用户ID列表
	userIDsStr := c.QueryArray("user_ids")
	if len(userIDsStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "必须提供至少一个用户ID"})
		return
	}

	// 解析用户ID
	userIDs := make([]uint64, 0, len(userIDsStr))
	for _, idStr := range userIDsStr {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID: " + idStr})
			return
		}
		userIDs = append(userIDs, id)
	}

	// 获取日期参数，默认为今天
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	// 解析日期
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日期格式，请使用YYYY-MM-DD格式"})
		return
	}

	// 获取统计数据
	stats, err := h.recordService.GetUsersDailyRecordStats(c, userIDs, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录统计失败: " + err.Error()})
		return
	}

	// 转换为数组形式返回
	result := make([]map[string]interface{}, 0, len(stats))
	for _, stat := range stats {
		// 格式化时间为字符串数组，方便前端处理
		recordTimeStrs := make([]string, len(stat.RecordTimes))
		for i, t := range stat.RecordTimes {
			recordTimeStrs[i] = t.Format("15:04:05")
		}

		result = append(result, map[string]interface{}{
			"user_id":      stat.UserID,
			"date":         stat.Date.Format("2006-01-02"),
			"count":        stat.Count,
			"total_time":   stat.TotalTime,
			"record_times": recordTimeStrs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": result,
		"date":  date.Format("2006-01-02"),
	})
}
