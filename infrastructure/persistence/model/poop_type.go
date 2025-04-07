package model

import (
	"record-project/domain/entity"
	"time"
)

// PoopType 屎的类型数据库模型
type PoopType struct {
	ID               uint64    `gorm:"primaryKey;column:id"`
	Name             string    `gorm:"type:varchar(50);not null;column:name;comment:类型名称"`
	Description      string    `gorm:"type:varchar(255);column:description;comment:类型描述"`
	Color            string    `gorm:"type:varchar(20);column:color;comment:颜色描述"`
	HealthIndication string    `gorm:"type:varchar(100);column:health_indication;comment:健康指示"`
	CreatedAt        time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime;column:updated_at;comment:更新时间"`
}

// TableName 指定表名
func (PoopType) TableName() string {
	return "poop_types"
}

// ToEntity 转换为领域实体
func (pt *PoopType) ToEntity() *entity.PoopType {
	return &entity.PoopType{
		ID:               pt.ID,
		Name:             pt.Name,
		Description:      pt.Description,
		Color:            pt.Color,
		HealthIndication: pt.HealthIndication,
		CreatedAt:        pt.CreatedAt,
		UpdatedAt:        pt.UpdatedAt,
	}
}

// FromEntity 从领域实体转换
func (pt *PoopType) FromEntity(poopType *entity.PoopType) {
	pt.ID = poopType.ID
	pt.Name = poopType.Name
	pt.Description = poopType.Description
	pt.Color = poopType.Color
	pt.HealthIndication = poopType.HealthIndication
	pt.CreatedAt = poopType.CreatedAt
	pt.UpdatedAt = poopType.UpdatedAt
}
