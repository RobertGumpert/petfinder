package repository

import (
	"advertservice/entity"
	"advertservice/mapper"
	"context"
	"gorm.io/gorm"
	"math"
	"time"
)

type GormSquareSearchModel struct {
	*gorm.DB
	compareTimeTruncate time.Duration
	areaSideLength      float64
	halfDiagonalSquare  float64
}

func NewGormSquareSearchModel(db *gorm.DB, compareTimeTruncate time.Duration, areaSideLength float64) *GormSquareSearchModel {
	s := &GormSquareSearchModel{
		DB :                 db,
		compareTimeTruncate: compareTimeTruncate,
		areaSideLength:      areaSideLength,
	}
	s.halfDiagonalSquare = s.getHalfDiagonalSquare(s.areaSideLength)
	return s
}

func (s *GormSquareSearchModel) FindAdverts(inputViewModel *mapper.FindAdvertsViewModel, ctx context.Context) ([]entity.Advert, error) {
	where := s.searchInAreaQuery(inputViewModel.GeoLatitude, inputViewModel.GeoLongitude)
	where.Limit(30)
	if inputViewModel.OnlyNotClosed {
		where.Not(map[string]interface{}{
			"date_close": nil,
		})
	}
	if inputViewModel.TypeLost && !inputViewModel.TypeAll {
		where.Where("ad_type = ? ", uint64(mapper.TypeLost))
	}
	if inputViewModel.TypeLost && !inputViewModel.TypeAll {
		where.Where("ad_type = ? ", uint64(mapper.TypeFound))
	}
	if !inputViewModel.AllOwners {
		if inputViewModel.OnlyOwnerAdverts {
			where.Where("ad_owner_id = ? ", inputViewModel.AdOwnerID)
		} else {
			where.Not("ad_owner_id = ? ", inputViewModel.AdOwnerID)
		}
	}
	if inputViewModel.Offset > 0 {
		where.Offset(int(inputViewModel.Offset))
	}
	var recordedListAdvertEntity []entity.Advert
	err := where.Find(&recordedListAdvertEntity).Error
	if err != nil {
		return nil, mapper.ErrorBadDataOperation
	}
	return recordedListAdvertEntity, nil
}
func (s *GormSquareSearchModel) SearchInArea(inputViewModel *mapper.SearchInAreaViewModel, ctx context.Context) ([]entity.Advert, error) {
	var (
		minLatitude, maxLatitude   = s.getMinMax(inputViewModel.GeoLatitude, s.halfDiagonalSquare)
		minLongitude, maxLongitude = s.getMinMax(inputViewModel.GeoLongitude, s.halfDiagonalSquare)
		result                     []entity.Advert
	)
	where := s.DB.Limit(30)
	where.Where(
		"(geo_latitude BETWEEN ? AND ?) AND (geo_longitude BETWEEN ? AND ?)",
		minLatitude,
		maxLatitude,
		minLongitude,
		maxLongitude,
	).Not(map[string]interface{}{
		"ad_owner_id": inputViewModel.AdOwnerID,
	})
	if inputViewModel.OnlyNotClosed {
		where.Not(map[string]interface{}{
			"date_close": nil,
		})
	}
	err := where.Find(&result)
	return result, err.Error
}

func (s *GormSquareSearchModel) searchInAreaQuery(lat, long float64) *gorm.DB {
	var (
		half                       = s.getHalfDiagonalSquare(s.areaSideLength)
		minLatitude, maxLatitude   = s.getMinMax(lat, half)
		minLongitude, maxLongitude = s.getMinMax(long, half)
	)
	where := s.DB
	where.Where(
		"(geo_latitude BETWEEN ? AND ?) AND (geo_longitude BETWEEN ? AND ?)",
		minLatitude,
		maxLatitude,
		minLongitude,
		maxLongitude,
	)
	return where
}

func (s *GormSquareSearchModel) getHalfDiagonalSquare(areaSideLength float64) float64 {
	return (math.Sqrt(2) * areaSideLength) / 2
}

func (s *GormSquareSearchModel) getMinMax(value, half float64) (min, max float64) {
	max = half + value
	min = value - half
	return min, max
}
