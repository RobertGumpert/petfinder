package service

import (
	"advertservice/entity"
	"advertservice/mapper"
	"advertservice/pckg/runtimeinfo"
	"advertservice/repository"
	"context"
	"log"
	"strings"
	"time"
)

type AdvertService struct {
	lifetimeOfFoundAnimalAdvert time.Duration
	lifetimeOfLostAnimalAdvert  time.Duration
	compareTimeTruncate         time.Duration
}

func NewAdvertService(lifetimeOfFoundAnimalAdvert time.Duration, lifetimeOfLostAnimalAdvert time.Duration, compareTimeTruncate time.Duration) *AdvertService {
	return &AdvertService{lifetimeOfFoundAnimalAdvert: lifetimeOfFoundAnimalAdvert, lifetimeOfLostAnimalAdvert: lifetimeOfLostAnimalAdvert, compareTimeTruncate: compareTimeTruncate}
}

func (s *AdvertService) CompareDates(comparedWith, verifiableWith time.Time, truncate time.Duration) (after, before, equal bool) {
	after, before, equal = false, false, false
	truncateCompared := comparedWith.Truncate(truncate)
	truncateVerifiable := verifiableWith.Truncate(truncate)
	after = truncateCompared.After(truncateVerifiable)
	before = truncateCompared.Before(truncateVerifiable)
	equal = truncateCompared.Equal(truncateVerifiable)
	return after, before, equal
}

func (s *AdvertService) GetLifetime(advertType mapper.TypeAdvert, dateCreate time.Time, truncate ...time.Duration) time.Time {
	var (
		dateClose time.Time
	)
	switch advertType {
	case mapper.TypeLost:
		dateClose = dateCreate.Add(s.lifetimeOfLostAnimalAdvert)
		break
	case mapper.TypeFound:
		dateClose = dateCreate.Add(s.lifetimeOfFoundAnimalAdvert)
		break
	default:
		dateClose = dateCreate
	}
	if len(truncate) == 0 {
		return dateClose
	}
	return dateClose.Truncate(truncate[0])
}

func (s *AdvertService) CreateAdvert(inputViewModel *mapper.CreateAdvertViewModel, db repository.AdvertRepository, ctx context.Context) (*mapper.AdvertViewModel, error) {
	err := inputViewModel.Validator()
	if err != nil {
		return nil, err
	}
	createAdvertEntity := inputViewModel.Mapper()
	if err = s.setLifetime(createAdvertEntity, false); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	if err = db.Create(createAdvertEntity, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	recordedAdvertEntity, err := db.EntityGet(createAdvertEntity, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.AdvertViewModel).Mapper(recordedAdvertEntity), nil
}

func (s *AdvertService) FindAdverts(inputViewModel *mapper.FindAdvertsViewModel, db repository.SearchModel, ctx context.Context) (*mapper.ListAdvertViewModel, error) {
	recordedListAdvertEntity, err := db.FindAdverts(inputViewModel, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return s.clusteringAndCloseExpire(recordedListAdvertEntity, nil, ctx, false, false), nil
}

func (s *AdvertService) SearchInArea(inputViewModel *mapper.SearchInAreaViewModel, db repository.SearchModel, ctx context.Context) (*mapper.ListAdvertViewModel, error) {
	recordedListAdvertEntity, err := db.SearchInArea(inputViewModel, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return s.clusteringAndCloseExpire(recordedListAdvertEntity, nil, ctx, false, false), nil
}

func (s *AdvertService) ListMyAdverts(inputViewModel *mapper.IdentifierOwnerViewModel, db repository.AdvertRepository, ctx context.Context) (*mapper.ListAdvertViewModel, error) {
	err := inputViewModel.Validator()
	if err != nil {
		return nil, err
	}
	var recordedListAdvertEntity []entity.Advert
	if recordedListAdvertEntity, err = db.EntityList(&entity.Advert{AdOwnerID: inputViewModel.AdOwnerID}, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return s.clusteringAndCloseExpire(recordedListAdvertEntity, db, ctx, true, true), nil
}

func (s *AdvertService) UpdateImage(inputViewModel *mapper.UpdateImageViewModel, db repository.AdvertRepository, ctx context.Context) error {
	err := inputViewModel.Validator()
	if err != nil {
		return err
	}
	advertId, mapNullableProperties := inputViewModel.Mapper()
	if err = db.MapUpdate(advertId, mapNullableProperties, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return mapper.ErrorBadDataOperation
	}
	return nil
}

func (s *AdvertService) Update(inputViewModel *mapper.UpdateAdvertViewModel, db repository.AdvertRepository, ctx context.Context) (*mapper.AdvertViewModel, error) {
	err := inputViewModel.Validator()
	if err != nil {
		return nil, err
	}
	updateAdvertEntity := inputViewModel.Mapper()
	recordedAdvertEntity, err := db.EntityGet(&entity.Advert{AdID: inputViewModel.AdID}, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	if inputViewModel.AdType != recordedAdvertEntity.AdType {
		if err = s.setLifetime(updateAdvertEntity, false); err != nil {
			go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return nil, mapper.ErrorNonValidData
		}
	}
	if err = db.EntityUpdate(updateAdvertEntity, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	if strings.TrimSpace(updateAdvertEntity.CommentText) != "" && recordedAdvertEntity.CommentText != updateAdvertEntity.CommentText {
		if err = db.MapUpdate(
			updateAdvertEntity.AdID,
			map[string]interface{}{
				"comment_text": updateAdvertEntity.CommentText,
			}, ctx); err != nil {
			go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return nil, mapper.ErrorBadDataOperation
		}
	}
	if updateAdvertEntity, err = db.EntityGet(updateAdvertEntity, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.AdvertViewModel).Mapper(updateAdvertEntity), nil
}

func (s *AdvertService) UpdateOwnerName(inputViewModel *mapper.IdentifierOwnerViewModel, db repository.AdvertRepository, ctx context.Context) error {
	err := inputViewModel.Validator()
	if err != nil {
		return err
	}
	var recordedListAdvertEntity []entity.Advert
	if recordedListAdvertEntity, err = db.EntityList(&entity.Advert{AdOwnerID: inputViewModel.AdOwnerID}, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return mapper.ErrorBadDataOperation
	}
	ids := make([]uint64, 0)
	mapUpdate := map[string]interface{}{
		"ad_owner_name": inputViewModel.AdOwnerName,
	}
	for _, advert := range recordedListAdvertEntity {
		ids = append(ids, advert.AdID)
	}
	if err := db.MapUpdateInID(ids, mapUpdate, ctx); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return err
	}
	return nil
}

func (s *AdvertService) RefreshAdvert(inputViewModel *mapper.IdentifierAdvertViewModel, db repository.AdvertRepository, ctx context.Context) (*mapper.UpdateLifetimeViewModel, error) {
	err := inputViewModel.Validator()
	if err != nil {
		return nil, err
	}
	refreshAdvertEntity := inputViewModel.Mapper()
	if refreshAdvertEntity, err = s.updateAdvertLifetime(refreshAdvertEntity, db, ctx, false); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.UpdateLifetimeViewModel).Mapper(refreshAdvertEntity), nil
}

func (s *AdvertService) CloseAdvert(inputViewModel *mapper.IdentifierAdvertViewModel, db repository.AdvertRepository, ctx context.Context) (*mapper.UpdateLifetimeViewModel, error) {
	err := inputViewModel.Validator()
	if err != nil {
		return nil, err
	}
	closeAdvertEntity := inputViewModel.Mapper()
	if closeAdvertEntity, err = s.updateAdvertLifetime(closeAdvertEntity, db, ctx, true); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.UpdateLifetimeViewModel).Mapper(closeAdvertEntity), nil
}

func (s *AdvertService) AdvertIsClosed(advertEntity *entity.Advert, t ...time.Duration) bool {
	if advertEntity.DateClose == nil {
		return true
	}
	var truncate time.Duration
	if len(t) != 0 {
		truncate = t[0]
	} else {
		truncate = s.compareTimeTruncate
	}
	now := time.Now().Truncate(truncate)
	after, _, _ := s.CompareDates(now, *advertEntity.DateClose, truncate)
	if after {
		return true
	}
	return false
}

func (s *AdvertService) clusteringAndCloseExpire(listAdvertEntity []entity.Advert, db repository.AdvertRepository, ctx context.Context, closeExpire, asyncCloseAdverts bool) *mapper.ListAdvertViewModel {
	var (
		lost         = make([]entity.Advert, 0)
		lostExpire   = make([]entity.Advert, 0)
		found        = make([]entity.Advert, 0)
		foundExpire  = make([]entity.Advert, 0)
		closedList   = make([]entity.Advert, 0)
		closeAdverts = func(closedList []entity.Advert, db repository.AdvertRepository, ctx context.Context) {
			if len(closedList) == 0 {
				return
			}
			ids := make([]uint64, 0)
			for _, advert := range closedList {
				ids = append(ids, advert.AdID)
			}
			if err := db.MapUpdateInID(ids, map[string]interface{}{
				"date_close": nil,
			}, ctx); err != nil {
				log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			}
		}
	)
	for _, advert := range listAdvertEntity {
		if advert.DateClose != nil {
			if s.AdvertIsClosed(&advert) {
				closedList = append(closedList, advert)
				advert.DateClose = nil
			}
		}
		switch advert.AdType {
		case uint64(mapper.TypeFound):
			if advert.DateClose == nil {
				foundExpire = append(foundExpire, advert)
			} else {
				found = append(found, advert)
			}
			break
		case uint64(mapper.TypeLost):
			if advert.DateClose == nil {
				lostExpire = append(lostExpire, advert)
			} else {
				lost = append(lost, advert)
			}
			break
		}
	}
	if closeExpire {
		if asyncCloseAdverts {
			go closeAdverts(closedList, db, ctx)
		} else {
			closeAdverts(closedList, db, ctx)
		}
	}
	return new(mapper.ListAdvertViewModel).Mapper(lost, lostExpire, found, foundExpire)
}

func (s *AdvertService) setLifetime(advertEntity *entity.Advert, closeAdvert bool) error {
	if !closeAdvert {
		dateCreate := time.Now()
		dateClose := s.GetLifetime(mapper.TypeAdvert(advertEntity.AdType), dateCreate)
		advertEntity.DateCreate = &dateCreate
		advertEntity.DateClose = &dateClose
	}
	if closeAdvert {
		advertEntity.DateClose = nil
	}
	return nil
}

func (s *AdvertService) updateAdvertLifetime(advertEntity *entity.Advert, db repository.AdvertRepository, ctx context.Context, closeAdvert bool) (*entity.Advert, error) {
	if !closeAdvert {
		if err := s.setLifetime(advertEntity, false); err != nil {
			go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return nil, err
		}
		return advertEntity, db.MapUpdate(advertEntity.AdID, map[string]interface{}{
			"date_create": advertEntity.DateCreate,
			"date_close":  advertEntity.DateClose,
		}, ctx)
	}
	if closeAdvert {
		if err := s.setLifetime(advertEntity, true); err != nil {
			go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return nil, err
		}
		return advertEntity, db.MapUpdate(advertEntity.AdID, map[string]interface{}{
			"date_close": nil,
		}, ctx)
	}
	return advertEntity, mapper.ErrorNonValidData
}
