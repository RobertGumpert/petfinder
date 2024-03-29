package repository

import (
	"advertservice/entity"
	"advertservice/mapper"
	"context"
)

type AdvertRepository interface {
	Create(advert *entity.Advert, ctx context.Context) error
	//
	EntityUpdate(advert *entity.Advert, ctx context.Context) error
	EntityGet(advert *entity.Advert, ctx context.Context) (*entity.Advert, error)
	EntityList(advert *entity.Advert, ctx context.Context) ([]entity.Advert, error)
	//
	MapUpdate(id uint64, fields map[string]interface{}, ctx context.Context) error
	MapUpdateInID(ids []uint64, fields map[string]interface{}, ctx context.Context) error
	MapGet(fields map[string]interface{}, ctx context.Context) (*entity.Advert, error)
	MapList(fields map[string]interface{}, ctx context.Context) ([]entity.Advert, error)
	//
	ListByID(id []uint64, ctx context.Context) ([]entity.Advert, error)
	//
	GetOrm() interface{}
}

type SearchModel interface {
	FindAdverts(inputViewModel *mapper.FindAdvertsViewModel, ctx context.Context) ([]entity.Advert, error)
	SearchInArea(inputViewModel *mapper.SearchInAreaViewModel, ctx context.Context) ([]entity.Advert, error)
}
