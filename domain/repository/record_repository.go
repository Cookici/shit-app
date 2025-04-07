package repository

import (
	"context"
	"record-project/domain/entity"
	"time"
)

// RecordRepository 记录仓储接口
type RecordRepository interface {
	FindByID(ctx context.Context, id uint64) (*entity.Record, error)
	FindByUserID(ctx context.Context, userID uint64, page, size int) ([]*entity.Record, int64, error)
	FindByDateRange(ctx context.Context, userID uint64, start, end time.Time, page, size int) ([]*entity.Record, error)
	Save(ctx context.Context, record *entity.Record) error
	Update(ctx context.Context, record *entity.Record) error
	Delete(ctx context.Context, id uint64) error
	CreateWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64, recordTagRepo RecordTagRepository) error
	UpdateWithTags(ctx context.Context, record *entity.Record, tagIDs []uint64, recordTagRepo RecordTagRepository) error
	CountRecordsByDateRange(ctx context.Context, userID uint64, start time.Time, end time.Time) (int64, error)
	GetGlobalRanking(ctx context.Context, start, end time.Time, limit int) ([]*entity.RankingItem, error)
	GetFriendRanking(ctx context.Context, userIDs []uint64, startDate, endDate time.Time, page, pageSize int) ([]*entity.RankingItem, int, error)
	GetUsersDailyRecordStats(ctx context.Context, userIDs []uint64, date time.Time) (map[uint64]*entity.DailyRecordStats, error)
}
