package service

import (
	"context"
	"record-project/domain/entity"
	"record-project/domain/repository"
)

// PoopTypeService 屎类型服务接口
type PoopTypeService interface {
	GetPoopTypeByID(ctx context.Context, id uint64) (*entity.PoopType, error)
	GetAllPoopTypes(ctx context.Context) ([]*entity.PoopType, error)
	GetPoopTypesByIDs(ctx context.Context, ids []uint64) ([]*entity.PoopType, error)
	CreatePoopType(ctx context.Context, poopType *entity.PoopType) error
	UpdatePoopType(ctx context.Context, poopType *entity.PoopType) error
	DeletePoopType(ctx context.Context, id uint64) error
}

// poopTypeService 屎的类型服务实现
type poopTypeService struct {
	poopTypeRepo repository.PoopTypeRepository
}

// NewPoopTypeService 创建屎的类型服务
func NewPoopTypeService(poopTypeRepo repository.PoopTypeRepository) PoopTypeService {
	return &poopTypeService{
		poopTypeRepo: poopTypeRepo,
	}
}

// GetPoopTypeByID 根据ID获取屎的类型
func (s *poopTypeService) GetPoopTypeByID(ctx context.Context, id uint64) (*entity.PoopType, error) {
	return s.poopTypeRepo.FindByID(ctx, id)
}

// GetAllPoopTypes 获取所有屎的类型
func (s *poopTypeService) GetAllPoopTypes(ctx context.Context) ([]*entity.PoopType, error) {
	return s.poopTypeRepo.FindAll(ctx)
}

// GetPoopTypesByIDs 批量获取多个屎类型
func (s *poopTypeService) GetPoopTypesByIDs(ctx context.Context, ids []uint64) ([]*entity.PoopType, error) {
	return s.poopTypeRepo.FindByIDs(ctx, ids)
}

// CreatePoopType 创建屎的类型
func (s *poopTypeService) CreatePoopType(ctx context.Context, poopType *entity.PoopType) error {
	return s.poopTypeRepo.Save(ctx, poopType)
}

// UpdatePoopType 更新屎的类型
func (s *poopTypeService) UpdatePoopType(ctx context.Context, poopType *entity.PoopType) error {
	return s.poopTypeRepo.Update(ctx, poopType)
}

// DeletePoopType 删除屎的类型
func (s *poopTypeService) DeletePoopType(ctx context.Context, id uint64) error {
	return s.poopTypeRepo.Delete(ctx, id)
}