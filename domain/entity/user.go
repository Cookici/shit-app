package entity

import "time"

// User 用户实体
type User struct {
	ID        uint64    `json:"id"`
	OpenID    string    `json:"open_id"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}