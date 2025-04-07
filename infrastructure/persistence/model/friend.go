package model

import (
	"record-project/domain/entity"
	"time"
)

// Friend 好友关系数据库模型
type Friend struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	UserID    uint64    `gorm:"index:idx_user_friend,unique;column:user_id;comment:用户ID"`
	FriendID  uint64    `gorm:"index:idx_user_friend,unique;column:friend_id;comment:好友ID"`
	Status    int8      `gorm:"type:tinyint;default:0;column:status;comment:关系状态: 0-待确认, 1-已确认, 2-已拒绝, 3-已删除"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;comment:更新时间"`
}

// TableName 指定表名
func (Friend) TableName() string {
	return "friends"
}

// ToEntity 转换为实体
func (f *Friend) ToEntity() *entity.Friend {
	return &entity.Friend{
		ID:        f.ID,
		UserID:    f.UserID,
		FriendID:  f.FriendID,
		Status:    f.Status,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

// FromEntity 从实体转换
func (f *Friend) FromEntity(friend *entity.Friend) {
	f.ID = friend.ID
	f.UserID = friend.UserID
	f.FriendID = friend.FriendID
	f.Status = friend.Status
	f.CreatedAt = friend.CreatedAt
	f.UpdatedAt = friend.UpdatedAt
}
