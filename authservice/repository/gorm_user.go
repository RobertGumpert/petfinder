package repository

import (
	"authservice/entity"
	"context"
	"gorm.io/gorm"
)



type UserGormRepository struct {
	*gorm.DB
}

func NewUserGormRepository(DB *gorm.DB) *UserGormRepository {
	return &UserGormRepository{DB: DB}
}

func (r *UserGormRepository) Create(model *entity.User, ctx context.Context) error {
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


func (r *UserGormRepository) EntityUpdate(model *entity.User, ctx context.Context) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entity.User{}).Where(&entity.User{UserID: model.UserID}).Updates(model).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *UserGormRepository) MapUpdate(id uint64, fields map[string]interface{}, ctx context.Context) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entity.User{}).Where(&entity.User{UserID: id}).Updates(fields).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//
func (r *UserGormRepository) GetByID(id uint64, ctx context.Context) (*entity.User, error) {
	var (
		result entity.User
	)
	err := r.DB.Where(&entity.User{UserID: id}).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserGormRepository) EntityGet(model *entity.User, ctx context.Context) (*entity.User, error) {
	var (
		result entity.User
	)
	err := r.DB.Where(model).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserGormRepository) MapGet(fields map[string]interface{}, ctx context.Context) (*entity.User, error) {
	var (
		result entity.User
	)
	err := r.DB.Where(fields).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//
func (r *UserGormRepository) ListByID(id []uint64, ctx context.Context) ([]entity.User, error) {
	var (
		result []entity.User
	)
	err := r.DB.Where("user_id IN ? ", id).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserGormRepository) EntityList(model *entity.User, ctx context.Context) ([]entity.User, error) {
	var (
		result []entity.User
	)
	err := r.DB.Where(model).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserGormRepository) MapList(fields map[string]interface{}, ctx context.Context) ([]entity.User, error) {
	var (
		result []entity.User
	)
	err := r.DB.Where(fields).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
