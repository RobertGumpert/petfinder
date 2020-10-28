package repository

import (
	"authservice/entity"
	"context"
)

type UserRepository interface {
	Create(user *entity.User, ctx context.Context) error
	//
	EntityUpdate(user *entity.User, ctx context.Context) error
	EntityGet(user *entity.User, ctx context.Context) (*entity.User, error)
	EntityList(user *entity.User, ctx context.Context) ([]entity.User, error)
	//
	MapUpdate(id uint64, fields map[string]interface{}, ctx context.Context) error
	MapGet(fields map[string]interface{}, ctx context.Context) (*entity.User, error)
	MapList(fields map[string]interface{}, ctx context.Context) ([]entity.User, error)
	//
	ListByID(id []uint64, ctx context.Context) ([]entity.User, error)
}
