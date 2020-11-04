package entity

import "time"

type Advert struct {
	//
	// Владелец
	//
	AdOwnerID   uint64 `gorm:"not null"`
	AdOwnerName string `gorm:"size:255;not null;"`
	//
	// Св-ва
	//
	AdID         uint64     `gorm:"primary_key;auto_increment"`
	DateCreate   *time.Time `gorm:"default:CURRENT_TIMESTAMP;not null;"`
	DateClose    *time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	AdType       uint64     `gorm:"size:10;not null;"`
	AnimalType   string     `gorm:"size:255;not null;"`
	AnimalBreed  string     `gorm:"size:255;not null;"`
	GeoLatitude  float64    `gorm:"size:10;not null;"`
	GeoLongitude float64    `gorm:"size:10;not null;"`
	CommentText  string     `gorm:"size:255;not null;"`
	ImageUrl     string     `gorm:"size:5000;"`
	//
	// gorm model
	//
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
}
