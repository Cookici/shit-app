package api

import (
	"net/http"
	"record-project/application/service"
	"record-project/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签API处理器
type TagHandler struct {
	tagService service.TagService
}

// NewTagHandler 创建标签API处理器
func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetTag 获取标签
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	tag, err := h.tagService.GetTagByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tag == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// GetAllTags 获取所有标签
func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var tag entity.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tagService.CreateTag(c, &tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var tag entity.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.ID = id
	if err := h.tagService.UpdateTag(c, &tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	if err := h.tagService.DeleteTag(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签删除成功"})
}
