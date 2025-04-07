package model

import (
	"record-project/domain/entity"
	"time"
)

// Tag 标签数据库模型
type Tag struct {
	ID          uint64    `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"type:varchar(50);not null;column:name;comment:标签名称"`
	Description string    `gorm:"type:varchar(255);column:description;comment:标签描述"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at;comment:更新时间"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// ToEntity 转换为领域实体
func (t *Tag) ToEntity() *entity.Tag {
	return &entity.Tag{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// FromEntity 从领域实体转换
func (t *Tag) FromEntity(tag *entity.Tag) {
	t.ID = tag.ID
	t.Name = tag.Name
	t.Description = tag.Description
	t.CreatedAt = tag.CreatedAt
	t.UpdatedAt = tag.UpdatedAt
}
