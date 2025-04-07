package repository

import (
	"context"
	"record-project/domain/entity"
)

// RecordTagRepository 记录标签关联仓储接口
type RecordTagRepository interface {
	SaveRecordTags(ctx context.Context, recordID uint64, tagIDs []uint64) error
	DeleteRecordTags(ctx context.Context, recordID uint64) error
	FindTagsByRecordID(ctx context.Context, recordID uint64) ([]*entity.Tag, error)
	FindTagsByRecordIDs(ctx context.Context, recordIDs []uint64) (map[uint64][]*entity.Tag, error)
	FindRecordsByTagID(ctx context.Context, tagID uint64) ([]*entity.Record, error)
}