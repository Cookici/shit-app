package repository

import (
	"context"
	"errors"
	"record-project/domain/entity"
	"record-project/domain/repository"
	"record-project/infrastructure/persistence/model"
	"time"

	"gorm.io/gorm"
)

// poopTypeRepository 屎的类型仓储实现
type poopTypeRepository struct {
	db *gorm.DB
}

// NewPoopTypeRepository 创建屎的类型仓储
func NewPoopTypeRepository(db *gorm.DB) repository.PoopTypeRepository {
	return &poopTypeRepository{db: db}
}

// FindByID 根据ID查找屎的类型
func (r *poopTypeRepository) FindByID(ctx context.Context, id uint64) (*entity.PoopType, error) {
	var poopTypeModel model.PoopType
	if err := r.db.WithContext(ctx).First(&poopTypeModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return poopTypeModel.ToEntity(), nil
}

// FindAll 查找所有屎的类型
func (r *poopTypeRepository) FindAll(ctx context.Context) ([]*entity.PoopType, error) {
	var poopTypeModels []model.PoopType

	if err := r.db.WithContext(ctx).Find(&poopTypeModels).Error; err != nil {
		return nil, err
	}

	poopTypes := make([]*entity.PoopType, len(poopTypeModels))
	for i, poopTypeModel := range poopTypeModels {
		poopTypes[i] = poopTypeModel.ToEntity()
	}

	return poopTypes, nil
}

// FindByIDs 根据ID列表查找多个屎类型
func (r *poopTypeRepository) FindByIDs(ctx context.Context, ids []uint64) ([]*entity.PoopType, error) {
	var poopTypeModels []model.PoopType
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&poopTypeModels).Error; err != nil {
		return nil, err
	}

	poopTypes := make([]*entity.PoopType, len(poopTypeModels))
	for i, ptm := range poopTypeModels {
		poopTypes[i] = ptm.ToEntity()
	}

	return poopTypes, nil
}

// Save 保存屎的类型
func (r *poopTypeRepository) Save(ctx context.Context, poopType *entity.PoopType) error {
	var poopTypeModel model.PoopType
	poopTypeModel.FromEntity(poopType)
	if err := r.db.WithContext(ctx).Create(&poopTypeModel).Error; err != nil {
		return err
	}
	poopType.ID = poopTypeModel.ID
	return nil
}

// Update 更新屎的类型
func (r *poopTypeRepository) Update(ctx context.Context, poopType *entity.PoopType) error {
	var poopTypeModel model.PoopType
	poopTypeModel.FromEntity(poopType)

	result := r.db.WithContext(ctx).Model(&model.PoopType{}).Where("id = ?", poopType.ID).Updates(map[string]interface{}{
		"name":        poopTypeModel.Name,
		"description": poopTypeModel.Description,
		"color":       poopTypeModel.Color,
		"updated_at":  time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete 删除屎的类型
func (r *poopTypeRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.PoopType{}, id).Error
}
