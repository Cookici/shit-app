package api

import (
	"net/http"
	"record-project/application/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PoopTypeHandler 屎的类型API处理器
type PoopTypeHandler struct {
	poopTypeService service.PoopTypeService
}

// NewPoopTypeHandler 创建屎的类型API处理器
func NewPoopTypeHandler(poopTypeService service.PoopTypeService) *PoopTypeHandler {
	return &PoopTypeHandler{
		poopTypeService: poopTypeService,
	}
}

// GetPoopType 获取屎的类型
func (h *PoopTypeHandler) GetPoopType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的屎的类型ID"})
		return
	}

	poopType, err := h.poopTypeService.GetPoopTypeByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if poopType == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "屎的类型不存在"})
		return
	}

	c.JSON(http.StatusOK, poopType)
}

// GetAllPoopTypes 获取所有屎的类型
func (h *PoopTypeHandler) GetAllPoopTypes(c *gin.Context) {
	poopTypes, err := h.poopTypeService.GetAllPoopTypes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, poopTypes)
}
