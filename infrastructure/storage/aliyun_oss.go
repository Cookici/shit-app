package storage

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"record-project/infrastructure/config"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// OSSService 阿里云OSS服务接口
type OSSService interface {
	UploadFile(data []byte, fileName string) (string, error)
	UploadFileWithReader(reader io.Reader, fileName string) (string, error)
}

// ossService 阿里云OSS服务实现
type ossService struct {
	client    *oss.Client
	bucket    *oss.Bucket
	urlPrefix string
}

// NewOSSService 创建阿里云OSS服务
func NewOSSService(cfg *config.AliyunConfig) (OSSService, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		return nil, err
	}

	return &ossService{
		client:    client,
		bucket:    bucket,
		urlPrefix: cfg.URLPrefix,
	}, nil
}

// UploadFile 上传文件
func (s *ossService) UploadFile(data []byte, fileName string) (string, error) {
	return s.UploadFileWithReader(bytes.NewReader(data), fileName)
}

// UploadFileWithReader 使用Reader上传文件
func (s *ossService) UploadFileWithReader(reader io.Reader, fileName string) (string, error) {
	// 生成唯一文件名
	ext := filepath.Ext(fileName)
	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期组织文件
	now := time.Now()
	objectKey := fmt.Sprintf("uploads/%d/%02d/%02d/%s", now.Year(), now.Month(), now.Day(), uniqueFileName)

	// 上传文件
	err := s.bucket.PutObject(objectKey, reader)
	if err != nil {
		return "", err
	}

	// 返回文件URL
	return s.urlPrefix + objectKey, nil
}
