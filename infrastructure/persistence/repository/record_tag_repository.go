package repository

import (
	"context"
	"record-project/domain/entity"
	"record-project/infrastructure/persistence/model"

	"gorm.io/gorm"
)

// RecordTagRepository 记录标签关联仓储接口
type RecordTagRepository interface {
	SaveRecordTags(ctx context.Context, recordID uint64, tagIDs []uint64) error
	DeleteRecordTags(ctx context.Context, recordID uint64) error
	FindTagsByRecordID(ctx context.Context, recordID uint64) ([]*entity.Tag, error)
	FindRecordsByTagID(ctx context.Context, tagID uint64) ([]*entity.Record, error)
	FindTagsByRecordIDs(ctx context.Context, recordIDs []uint64) (map[uint64][]*entity.Tag, error)
}

// recordTagRepository 记录标签关联仓储实现
type recordTagRepository struct {
	db *gorm.DB
}

// NewRecordTagRepository 创建记录标签关联仓储
func NewRecordTagRepository(db *gorm.DB) RecordTagRepository {
	return &recordTagRepository{db: db}
}

// SaveRecordTags 保存记录的标签关联
func (r *recordTagRepository) SaveRecordTags(ctx context.Context, recordID uint64, tagIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除原有关联
		if err := tx.Where("record_id = ?", recordID).Delete(&model.RecordTag{}).Error; err != nil {
			return err
		}

		// 如果没有新标签，直接返回
		if len(tagIDs) == 0 {
			return nil
		}

		// 准备批量插入的数据
		recordTags := make([]model.RecordTag, len(tagIDs))
		for i, tagID := range tagIDs {
			recordTags[i] = model.RecordTag{
				RecordID: recordID,
				TagID:    tagID,
			}
		}

		// 批量插入新关联
		if err := tx.Create(&recordTags).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteRecordTags 删除记录的标签关联
func (r *recordTagRepository) DeleteRecordTags(ctx context.Context, recordID uint64) error {
	return r.db.WithContext(ctx).Where("record_id = ?", recordID).Delete(&model.RecordTag{}).Error
}

// FindTagsByRecordID 根据记录ID查找标签
func (r *recordTagRepository) FindTagsByRecordID(ctx context.Context, recordID uint64) ([]*entity.Tag, error) {
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

// FindRecordsByTagID 根据标签ID查找记录
func (r *recordTagRepository) FindRecordsByTagID(ctx context.Context, tagID uint64) ([]*entity.Record, error) {
	var recordTags []model.RecordTag
	if err := r.db.WithContext(ctx).Where("tag_id = ?", tagID).Find(&recordTags).Error; err != nil {
		return nil, err
	}

	if len(recordTags) == 0 {
		return []*entity.Record{}, nil
	}

	var recordIDs []uint64
	for _, rt := range recordTags {
		recordIDs = append(recordIDs, rt.RecordID)
	}

	var recordModels []model.Record
	if err := r.db.WithContext(ctx).Where("id IN ?", recordIDs).Find(&recordModels).Error; err != nil {
		return nil, err
	}

	records := make([]*entity.Record, len(recordModels))
	for i, recordModel := range recordModels {
		records[i] = recordModel.ToEntity()
	}

	return records, nil
}

// FindTagsByRecordIDs 批量查找多个记录的标签
func (r *recordTagRepository) FindTagsByRecordIDs(ctx context.Context, recordIDs []uint64) (map[uint64][]*entity.Tag, error) {
	var recordTags []model.RecordTag
	if err := r.db.WithContext(ctx).Where("record_id IN ?", recordIDs).Find(&recordTags).Error; err != nil {
		return nil, err
	}

	// 如果没有关联记录，返回空映射
	if len(recordTags) == 0 {
		return make(map[uint64][]*entity.Tag), nil
	}

	// 收集所有标签ID
	var tagIDs []uint64
	recordToTagIDs := make(map[uint64][]uint64)

	for _, rt := range recordTags {
		tagIDs = append(tagIDs, rt.TagID)
		recordToTagIDs[rt.RecordID] = append(recordToTagIDs[rt.RecordID], rt.TagID)
	}

	// 批量查询所有标签
	var tagModels []model.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", tagIDs).Find(&tagModels).Error; err != nil {
		return nil, err
	}

	// 构建标签ID到标签实体的映射
	tagMap := make(map[uint64]*entity.Tag)
	for _, tm := range tagModels {
		tag := tm.ToEntity()
		tagMap[tm.ID] = tag
	}

	// 构建记录ID到标签列表的映射
	result := make(map[uint64][]*entity.Tag)
	for recordID, tagIDs := range recordToTagIDs {
		tags := make([]*entity.Tag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			if tag, exists := tagMap[tagID]; exists {
				tags = append(tags, tag)
			}
		}
		result[recordID] = tags
	}

	return result, nil
}
