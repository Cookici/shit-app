package model

import (
	"record-project/domain/entity"
	"time"
)

// Record 拉屎记录数据库模型
type Record struct {
	ID         uint64    `gorm:"primaryKey;column:id"`
	UserID     uint64    `gorm:"not null;index;column:user_id;comment:用户ID"`
	RecordTime time.Time `gorm:"not null;column:record_time;comment:拉屎时间"`
	Duration   int       `gorm:"column:duration;comment:持续时间(秒)"`
	PoopTypeID uint64    `gorm:"column:poop_type_id;comment:屎的类型ID"`
	Note       string    `gorm:"type:text;column:note;comment:备注"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at;comment:更新时间"`
}

// TableName 指定表名
func (Record) TableName() string {
	return "records"
}

// ToEntity 转换为领域实体
func (r *Record) ToEntity() *entity.Record {
	return &entity.Record{
		ID:         r.ID,
		UserID:     r.UserID,
		RecordTime: r.RecordTime,
		Duration:   r.Duration,
		PoopTypeID: r.PoopTypeID,
		Note:       r.Note,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// FromEntity 从领域实体转换
func (r *Record) FromEntity(record *entity.Record) {
	r.ID = record.ID
	r.UserID = record.UserID
	r.RecordTime = record.RecordTime
	r.Duration = record.Duration
	r.PoopTypeID = record.PoopTypeID
	r.Note = record.Note
	r.CreatedAt = record.CreatedAt
	r.UpdatedAt = record.UpdatedAt
}
