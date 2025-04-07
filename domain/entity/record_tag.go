package entity

import "time"

// RecordTag 记录和标签的关联实体
type RecordTag struct {
	ID        uint64    `json:"id"`
	RecordID  uint64    `json:"record_id"`
	TagID     uint64    `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}