package repository

import (
	"advertservice/entity"
	"context"
	"gorm.io/gorm"
)

type GormAdvertRepository struct {
	*gorm.DB
}

func NewGormAdvertRepository(DB *gorm.DB) *GormAdvertRepository {
	return &GormAdvertRepository{DB: DB}
}

func (r *GormAdvertRepository) GetOrm() interface{} {
	return r.DB
}

func (r *GormAdvertRepository) Create(model *entity.Advert, ctx context.Context) error {
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


func (r *GormAdvertRepository) EntityUpdate(model *entity.Advert, ctx context.Context) error {
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

func (r *GormAdvertRepository) MapUpdate(id uint64, fields map[string]interface{}, ctx context.Context) error {
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
func (r *GormAdvertRepository) GetByID(id uint64, ctx context.Context) (*entity.Advert, error) {
	var (
		result entity.Advert
	)
	err := r.DB.Where(&entity.Advert{AdID: id}).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *GormAdvertRepository) EntityGet(model *entity.Advert, ctx context.Context) (*entity.Advert, error) {
	var (
		result entity.Advert
	)
	err := r.DB.Where(model).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *GormAdvertRepository) MapGet(fields map[string]interface{}, ctx context.Context) (*entity.Advert, error) {
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
func (r *GormAdvertRepository) ListByID(id []uint64, ctx context.Context) ([]entity.Advert, error) {
	var (
		result []entity.Advert
	)
	err := r.DB.Where("user_id IN ? ", id).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GormAdvertRepository) EntityList(model *entity.Advert, ctx context.Context) ([]entity.Advert, error) {
	var (
		result []entity.Advert
	)
	err := r.DB.Where(model).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GormAdvertRepository) MapList(fields map[string]interface{}, ctx context.Context) ([]entity.Advert, error) {
	var (
		result []entity.Advert
	)
	err := r.DB.Where(fields).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GormAdvertRepository) MapUpdateInID(ids []uint64, fields map[string]interface{}, ctx context.Context) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entity.Advert{}).Where("ad_id IN ? ", ids).Updates(fields).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}