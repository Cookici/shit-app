package repository

import (
	"context"
	"errors"
	"record-project/domain/entity"
	"record-project/domain/repository"
	"record-project/infrastructure/persistence/model"
	"time"

	"gorm.io/gorm"
)

// tagRepository 标签仓储实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

// FindByID 根据ID查找标签
func (r *tagRepository) FindByID(ctx context.Context, id uint64) (*entity.Tag, error) {
	var tagModel model.Tag
	if err := r.db.WithContext(ctx).First(&tagModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tagModel.ToEntity(), nil
}

// FindAll 查找所有标签
func (r *tagRepository) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	var tagModels []model.Tag

	if err := r.db.WithContext(ctx).Find(&tagModels).Error; err != nil {
		return nil, err
	}

	tags := make([]*entity.Tag, len(tagModels))
	for i, tagModel := range tagModels {
		tags[i] = tagModel.ToEntity()
	}

	return tags, nil
}

// FindByRecordID 根据记录ID查找标签
func (r *tagRepository) FindByRecordID(ctx context.Context, recordID uint64) ([]*entity.Tag, error) {
	var recordTags []model.RecordTag
	if err := r.db.WithContext(ctx).Where("record_id = ?", recordID).Find(&recordTags).Error; err != nil {
		return nil, err
	}

	if len(recordTags) == 0 {
		return []*entity.Tag{}, nil
	}

	var tagIDs []uint64
	for _, rt := range recordTags {
		tagIDs = append(tagIDs, rt.TagID)
	}

	var tagModels []model.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", tagIDs).Find(&tagModels).Error; err != nil {
		return nil, err
	}

	tags := make([]*entity.Tag, len(tagModels))
	for i, tagModel := range tagModels {
		tags[i] = tagModel.ToEntity()
	}

	return tags, nil
}

// Save 保存标签
func (r *tagRepository) Save(ctx context.Context, tag *entity.Tag) error {
	var tagModel model.Tag
	tagModel.FromEntity(tag)
	if err := r.db.WithContext(ctx).Create(&tagModel).Error; err != nil {
		return err
	}
	tag.ID = tagModel.ID
	return nil
}

// Update 更新标签
func (r *tagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	var tagModel model.Tag
	tagModel.FromEntity(tag)

	result := r.db.WithContext(ctx).Model(&model.Tag{}).Where("id = ?", tag.ID).Updates(map[string]interface{}{
		"name":        tagModel.Name,
		"description": tagModel.Description,
		"updated_at":  time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete 删除标签
func (r *tagRepository) Delete(ctx context.Context, id uint64) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除关联的记录标签关系
		if err := tx.Where("tag_id = ?", id).Delete(&model.RecordTag{}).Error; err != nil {
			return err
		}

		// 再删除标签本身
		return tx.Delete(&model.Tag{}, id).Error
	})
}
