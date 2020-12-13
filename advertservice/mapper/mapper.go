package mapper

import (
	"advertservice/entity"
	"strings"
	"time"
)

type UserViewModel struct {
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func (m *UserViewModel) Validator() error {
	if m.UserID == 0 || strings.TrimSpace(m.Name) == "" {
		return ErrorNonValidData
	}
	return nil
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type AdvertViewModel struct {
	AdOwnerID    uint64     `json:"ad_owner_id"`
	AdOwnerName  string     `json:"ad_owner_name"`
	AdType       uint64     `json:"ad_type"`
	AdID         uint64     `json:"ad_id"`
	AnimalType   string     `json:"animal_type"`
	AnimalBreed  string     `json:"animal_breed"`
	GeoLatitude  float64    `json:"geo_latitude"`
	GeoLongitude float64    `json:"geo_longitude"`
	CommentText  string     `json:"comment_text"`
	ImageUrl     string     `json:"image_url"`
	DateCreate   *time.Time `json:"date_create"`
	DateClose    *time.Time `json:"date_close"`
}

func (m *AdvertViewModel) Mapper(advert *entity.Advert) *AdvertViewModel {
	m.AdID = advert.AdID
	m.AdOwnerID = advert.AdOwnerID
	m.AdOwnerName = advert.AdOwnerName
	m.AdType = advert.AdType
	m.AnimalBreed = advert.AnimalBreed
	m.AnimalType = advert.AnimalType
	m.GeoLatitude = advert.GeoLatitude
	m.GeoLongitude = advert.GeoLongitude
	m.CommentText = advert.CommentText
	m.ImageUrl = advert.ImageUrl
	m.DateCreate = advert.DateCreate
	m.DateClose = advert.DateClose
	return m
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type LostListViewModel struct {
	List   []*AdvertViewModel `json:"list"`
	Expire []*AdvertViewModel `json:"expire"`
}

func (m *LostListViewModel) Mapper(list, expire []entity.Advert) *LostListViewModel {
	for _, advert := range list {
		m.List = append(m.List, new(AdvertViewModel).Mapper(&advert))
	}
	for _, advert := range expire {
		m.Expire = append(m.Expire, new(AdvertViewModel).Mapper(&advert))
	}
	return m
}

type FoundListViewModel struct {
	List   []*AdvertViewModel `json:"list"`
	Expire []*AdvertViewModel `json:"expire"`
}

func (m *FoundListViewModel) Mapper(list, expire []entity.Advert) *FoundListViewModel {
	for _, advert := range list {
		m.List = append(m.List, new(AdvertViewModel).Mapper(&advert))
	}
	for _, advert := range expire {
		m.Expire = append(m.Expire, new(AdvertViewModel).Mapper(&advert))
	}
	return m
}

type ListAdvertViewModel struct {
	Lost  *LostListViewModel  `json:"lost"`
	Found *FoundListViewModel `json:"found"`
}

func (m *ListAdvertViewModel) Mapper(lost, lostExpire, found, foundExpire []entity.Advert) *ListAdvertViewModel {
	m.Lost = &LostListViewModel{
		List:   make([]*AdvertViewModel, 0),
		Expire: make([]*AdvertViewModel, 0),
	}
	m.Found = &FoundListViewModel{
		List:   make([]*AdvertViewModel, 0),
		Expire: make([]*AdvertViewModel, 0),
	}
	m.Lost.Mapper(lost, lostExpire)
	m.Found.Mapper(found, foundExpire)
	return m
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type FindAdvertsViewModel struct {
	OnlyNotClosed bool `json:"only_not_closed"`
	//
	TypeAll   bool `json:"type_all"`
	TypeFound bool `json:"type_found"`
	TypeLost  bool `json:"type_lost"`
	//
	AllOwners        bool   `json:"all_owners"`
	OnlyOwnerAdverts bool   `json:"only_owner_adverts"`
	NotOwnerAdverts  bool   `json:"not_owner_adverts"`
	AdOwnerID        uint64 `json:"ad_owner_id"`
	//
	SearchRadius float64 `json:"search_radius"`
	GeoLatitude  float64 `json:"geo_latitude"`
	GeoLongitude float64 `json:"geo_longitude"`
	//
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type SearchInAreaViewModel struct {
	AdOwnerID     uint64  `json:"ad_owner_id"`
	OnlyNotClosed bool    `json:"only_not_closed"`
	GeoLongitude  float64 `json:"geo_longitude"`
	GeoLatitude   float64 `json:"geo_latitude"`
}

func (m *SearchInAreaViewModel) Validator() error {
	if m.AdOwnerID == 0 {
		return ErrorNonValidData
	}
	return nil
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type CreateAdvertViewModel struct {
	AdOwnerID    uint64  `json:"ad_owner_id"`
	AdOwnerName  string  `json:"ad_owner_name"`
	AdType       uint64  `json:"ad_type"`
	AnimalType   string  `json:"animal_type"`
	AnimalBreed  string  `json:"animal_breed"`
	GeoLatitude  float64 `json:"geo_latitude"`
	GeoLongitude float64 `json:"geo_longitude"`
	CommentText  string  `json:"comment_text"`
}

func (m *CreateAdvertViewModel) Validator() error {
	if m.AdOwnerID == 0 || strings.TrimSpace(m.AdOwnerName) == "" {
		return ErrorNonValidData
	}
	if strings.TrimSpace(m.CommentText) == "" {
		m.CommentText = " "
	}
	if m.AdType != uint64(TypeLost) && m.AdType != uint64(TypeFound) {
		return ErrorNonValidData
	}
	if m.GeoLatitude > MaxLatitude || m.GeoLatitude < MinLatitude {
		return ErrorNonValidData
	}
	if m.GeoLongitude > MaxLongitude || m.GeoLongitude < MinLongitude {
		return ErrorNonValidData
	}
	return nil
}

func (m *CreateAdvertViewModel) Mapper() *entity.Advert {
	return &entity.Advert{
		AdOwnerID:    m.AdOwnerID,
		AdOwnerName:  m.AdOwnerName,
		AdType:       m.AdType,
		AnimalType:   m.AnimalType,
		AnimalBreed:  m.AnimalBreed,
		GeoLatitude:  m.GeoLatitude,
		GeoLongitude: m.GeoLatitude,
		CommentText:  m.CommentText,
	}
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type UpdateImageViewModel struct {
	AdID     uint64 `json:"ad_id"`
	ImageUrl string `json:"image_url"`
}

func (m *UpdateImageViewModel) Validator() error {
	if m.AdID == 0 || strings.TrimSpace(m.ImageUrl) == "" {
		return ErrorNonValidData
	}
	return nil
}

func (m *UpdateImageViewModel) Mapper() (uint64, map[string]interface{}) {
	return m.AdID, map[string]interface{}{
		"image_url": m.ImageUrl,
	}
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type UpdateAdvertViewModel struct {
	AdID uint64 `json:"ad_id"`
	//
	AdOwnerName  string  `json:"ad_owner_name"`
	AdType       uint64  `json:"ad_type"`
	AnimalType   string  `json:"animal_type"`
	AnimalBreed  string  `json:"animal_breed"`
	GeoLatitude  float64 `json:"geo_latitude"`
	GeoLongitude float64 `json:"geo_longitude"`
	CommentText  string  `json:"comment_text"`
}

func (m *UpdateAdvertViewModel) Validator() error {
	if m.AdID == 0 {
		return ErrorNonValidData
	}
	if m.GeoLatitude > MaxLatitude || m.GeoLatitude < MinLatitude {
		return ErrorNonValidData
	}
	if m.GeoLongitude > MaxLongitude || m.GeoLongitude < MinLongitude {
		return ErrorNonValidData
	}
	return nil
}

func (m *UpdateAdvertViewModel) Mapper() *entity.Advert {
	return &entity.Advert{
		AdOwnerName:  m.AdOwnerName,
		AdID:         m.AdID,
		AdType:       m.AdType,
		AnimalType:   m.AnimalType,
		AnimalBreed:  m.AnimalBreed,
		GeoLatitude:  m.GeoLatitude,
		GeoLongitude: m.GeoLongitude,
	}
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type IdentifierAdvertViewModel struct {
	AdID   uint64 `json:"ad_id"`
	AdType uint64 `json:"ad_type"`
}

func (m *IdentifierAdvertViewModel) Validator() error {
	if m.AdID == 0 {
		return ErrorNonValidData
	}
	if m.AdType != uint64(TypeLost) && m.AdType != uint64(TypeFound) {
		return ErrorNonValidData
	}
	return nil
}

func (m *IdentifierAdvertViewModel) Mapper() *entity.Advert {
	return &entity.Advert{
		AdID:   m.AdID,
		AdType: m.AdType,
	}
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type IdentifierOwnerViewModel struct {
	AdOwnerID   uint64 `json:"ad_owner_id"`
	AdOwnerName string `json:"ad_owner_name"`
}

func (m *IdentifierOwnerViewModel) Validator() error {
	if m.AdOwnerID == 0 {
		return ErrorNonValidData
	}
	if m.AdOwnerID == 0 && strings.TrimSpace(m.AdOwnerName) == "" {
		return ErrorNonValidData
	}
	return nil
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type UpdateLifetimeViewModel struct {
	AdID       uint64     `json:"ad_id"`
	DateCreate *time.Time `json:"date_create"`
	DateClose  *time.Time `json:"date_close"`
}

func (m *UpdateLifetimeViewModel) Mapper(advert *entity.Advert) *UpdateLifetimeViewModel {
	m.AdID = advert.AdID
	m.DateCreate = advert.DateCreate
	m.DateClose = advert.DateClose
	return m
}

//
//----------------------------------------------------------------------------------------------------------------------
//
