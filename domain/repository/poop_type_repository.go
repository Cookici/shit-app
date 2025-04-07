package repository

import (
	"context"
	"record-project/domain/entity"
)

// PoopTypeRepository 屎类型仓储接口
type PoopTypeRepository interface {
	FindByID(ctx context.Context, id uint64) (*entity.PoopType, error)
	FindAll(ctx context.Context) ([]*entity.PoopType, error)
	Save(ctx context.Context, poopType *entity.PoopType) error
	Update(ctx context.Context, poopType *entity.PoopType) error
	Delete(ctx context.Context, id uint64) error
    
    // FindByIDs 根据ID列表查找多个屎类型
    FindByIDs(ctx context.Context, ids []uint64) ([]*entity.PoopType, error)
}
