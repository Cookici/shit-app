package repository

import (
	"context"
	"record-project/domain/entity"
)

// TagRepository 标签仓储接口
type TagRepository interface {
	FindByID(ctx context.Context, id uint64) (*entity.Tag, error)
	FindAll(ctx context.Context) ([]*entity.Tag, error)
	FindByRecordID(ctx context.Context, recordID uint64) ([]*entity.Tag, error)
	Save(ctx context.Context, tag *entity.Tag) error
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uint64) error
}
