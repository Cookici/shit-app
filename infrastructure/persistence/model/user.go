package model

import (
	"record-project/domain/entity"
	"time"
)

// User 用户数据库模型
type User struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	OpenID    string    `gorm:"uniqueIndex;not null;type:varchar(100);column:open_id;comment:微信用户唯一标识"`
	Nickname  string    `gorm:"type:varchar(50);default:用户;column:nickname;comment:用户昵称"`
	AvatarURL string    `gorm:"type:varchar(255);column:avatar_url;comment:用户头像URL"`
	Status    int8      `gorm:"type:tinyint;default:1;column:status;comment:用户状态: 1-正常, 0-禁用"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;comment:更新时间"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// ToEntity 转换为领域实体
func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        u.ID,
		OpenID:    u.OpenID,
		Nickname:  u.Nickname,
		AvatarURL: u.AvatarURL,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromEntity 从领域实体转换
func (u *User) FromEntity(user *entity.User) {
	u.ID = user.ID
	u.OpenID = user.OpenID
	u.Nickname = user.Nickname
	u.AvatarURL = user.AvatarURL
	u.Status = user.Status
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}
