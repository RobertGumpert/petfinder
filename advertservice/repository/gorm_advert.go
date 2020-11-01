package repository

import (
	"advertservice/entity"
	"context"
	"gorm.io/gorm"
)

type AdvertGormRepository struct {
	*gorm.DB
}

func NewAdvertGormRepository(DB *gorm.DB) *AdvertGormRepository {
	return &AdvertGormRepository{DB: DB}
}

func (r *AdvertGormRepository) Create(model *entity.Advert, ctx context.Context) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(model).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}


func (r *AdvertGormRepository) EntityUpdate(model *entity.Advert, ctx context.Context) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entity.Advert{}).Where(&entity.Advert{AdID: model.AdID}).Updates(model).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *AdvertGormRepository) MapUpdate(id uint64, fields map[string]interface{}, ctx context.Context) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entity.Advert{}).Where(&entity.Advert{AdID: id}).Updates(fields).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//
func (r *AdvertGormRepository) GetByID(id uint64, ctx context.Context) (*entity.Advert, error) {
	var (
		result entity.Advert
	)
	err := r.DB.Where(&entity.Advert{AdID: id}).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AdvertGormRepository) EntityGet(model *entity.Advert, ctx context.Context) (*entity.Advert, error) {
	var (
		result entity.Advert
	)
	err := r.DB.Where(model).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AdvertGormRepository) MapGet(fields map[string]interface{}, ctx context.Context) (*entity.Advert, error) {
	var (
		result entity.Advert
	)
	err := r.DB.Where(fields).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//
func (r *AdvertGormRepository) ListByID(id []uint64, ctx context.Context) ([]entity.Advert, error) {
	var (
		result []entity.Advert
	)
	err := r.DB.Where("user_id IN ? ", id).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AdvertGormRepository) EntityList(model *entity.Advert, ctx context.Context) ([]entity.Advert, error) {
	var (
		result []entity.Advert
	)
	err := r.DB.Where(model).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AdvertGormRepository) MapList(fields map[string]interface{}, ctx context.Context) ([]entity.Advert, error) {
	var (
		result []entity.Advert
	)
	err := r.DB.Where(fields).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
