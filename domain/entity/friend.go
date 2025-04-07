package entity

import "time"

// Friend 好友关系实体
type Friend struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	FriendID   uint64    `json:"friend_id"`
	Status     int8      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FriendUser *User     `json:"friend_user,omitempty"` // 好友的用户信息
}