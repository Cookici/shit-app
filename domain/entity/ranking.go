package entity

// RankingItem 排行榜项目
type RankingItem struct {
	Rank       uint64 `json:"rank"`
	UserID     uint64 `json:"user_id"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	RecordCount int64  `json:"record_count"`
	TotalDuration int64 `json:"total_duration"`
}