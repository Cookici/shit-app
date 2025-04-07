package entity

import "time"

// DailyRecordStats 用户每日拉屎记录统计
type DailyRecordStats struct {
	UserID      uint64    `json:"user_id"`
	Date        time.Time `json:"date"`
	Count       int       `json:"count"`         // 记录总数
	TotalTime   int       `json:"total_time"`    // 总时长(秒)
	RecordTimes []time.Time `json:"record_times"` // 记录时间列表
}