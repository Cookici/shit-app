package model

import (
	"record-project/domain/entity"
	"time"
)

// RecordTag 记录和标签的关联数据库模型
type RecordTag struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	RecordID  uint64    `gorm:"not null;uniqueIndex:idx_record_tag;column:record_id;comment:记录ID"`
	TagID     uint64    `gorm:"not null;uniqueIndex:idx_record_tag;column:tag_id;comment:标签ID"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
}

// TableName 指定表名
func (RecordTag) TableName() string {
	return "record_tags"
}

// ToEntity 转换为领域实体
func (rt *RecordTag) ToEntity() *entity.RecordTag {
	return &entity.RecordTag{
		ID:        rt.ID,
		RecordID:  rt.RecordID,
		TagID:     rt.TagID,
		CreatedAt: rt.CreatedAt,
	}
}

// FromEntity 从领域实体转换
func (rt *RecordTag) FromEntity(recordTag *entity.RecordTag) {
	rt.ID = recordTag.ID
	rt.RecordID = recordTag.RecordID
	rt.TagID = recordTag.TagID
	rt.CreatedAt = recordTag.CreatedAt
}
