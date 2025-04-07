package entity

import "time"

// PoopType 屎的类型实体
type PoopType struct {
	ID               uint64    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Color            string    `json:"color"`
	HealthIndication string    `json:"health_indication"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}