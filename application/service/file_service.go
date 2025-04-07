package service

import (
	"context"
	"record-project/infrastructure/storage"
)

// FileService 文件服务接口
type FileService interface {
	UploadFile(ctx context.Context, data []byte, fileName string) (string, error)
}

// fileService 文件服务实现
type fileService struct {
	ossService storage.OSSService
}

// NewFileService 创建文件服务
func NewFileService(ossService storage.OSSService) FileService {
	return &fileService{
		ossService: ossService,
	}
}

// UploadFile 上传文件
func (s *fileService) UploadFile(ctx context.Context, data []byte, fileName string) (string, error) {
	return s.ossService.UploadFile(data, fileName)
}