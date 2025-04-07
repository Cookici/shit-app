package service

import (
	"context"
	"record-project/domain/entity"
	"record-project/domain/repository"
)

// TagService 标签服务接口
type TagService interface {
	GetTagByID(ctx context.Context, id uint64) (*entity.Tag, error)
	GetAllTags(ctx context.Context) ([]*entity.Tag, error)
	GetTagsByRecordIDs(ctx context.Context, recordIDs []uint64) (map[uint64][]*entity.Tag, error)
	CreateTag(ctx context.Context, tag *entity.Tag) error
	UpdateTag(ctx context.Context, tag *entity.Tag) error
	DeleteTag(ctx context.Context, id uint64) error
}

// tagService 标签服务实现
type tagService struct {
	tagRepo       repository.TagRepository
	recordTagRepo repository.RecordTagRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo repository.TagRepository, recordTagRepo repository.RecordTagRepository) TagService {
	return &tagService{
		tagRepo:       tagRepo,
		recordTagRepo: recordTagRepo,
	}
}

// GetTagByID 根据ID获取标签
func (s *tagService) GetTagByID(ctx context.Context, id uint64) (*entity.Tag, error) {
	return s.tagRepo.FindByID(ctx, id)
}

// GetAllTags 获取所有标签
func (s *tagService) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	return s.tagRepo.FindAll(ctx)
}

// GetTagsByRecordIDs 批量获取多个记录的标签
func (s *tagService) GetTagsByRecordIDs(ctx context.Context, recordIDs []uint64) (map[uint64][]*entity.Tag, error) {
	return s.recordTagRepo.FindTagsByRecordIDs(ctx, recordIDs)
}

// CreateTag 创建标签
func (s *tagService) CreateTag(ctx context.Context, tag *entity.Tag) error {
	return s.tagRepo.Save(ctx, tag)
}

// UpdateTag 更新标签
func (s *tagService) UpdateTag(ctx context.Context, tag *entity.Tag) error {
	return s.tagRepo.Update(ctx, tag)
}

// DeleteTag 删除标签
func (s *tagService) DeleteTag(ctx context.Context, id uint64) error {
	return s.tagRepo.Delete(ctx, id)
}