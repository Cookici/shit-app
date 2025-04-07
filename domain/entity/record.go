package entity

import "time"

// Record 拉屎记录实体
type Record struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	RecordTime time.Time `json:"record_time"`
	Duration   int       `json:"duration"`
	PoopTypeID uint64    `json:"poop_type_id"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联对象，不存储在数据库中
	User     *User     `json:"user,omitempty" gorm:"-"`
	PoopType *PoopType `json:"poop_type,omitempty" gorm:"-"`
	Tags     []*Tag    `json:"tags,omitempty" gorm:"-"`
}
