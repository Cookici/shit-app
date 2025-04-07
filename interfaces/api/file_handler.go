package api

import (
	"bytes"
	"io"
	"net/http"
	"record-project/application/service"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件API处理器
type FileHandler struct {
	fileService service.FileService
}

// NewFileHandler 创建文件API处理器
func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	// 读取文件内容
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	// 上传文件
	url, err := h.fileService.UploadFile(c, buffer.Bytes(), header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"message": "文件上传成功",
	})
}
