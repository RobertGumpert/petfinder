package entity

import "time"

type User struct {
	UserID       uint64 `gorm:"primary_key;auto_increment"`
	Telephone    string `gorm:"size:16;not null;unique"`
	Password     string `gorm:"size:255;not null;"`
	Email        string `gorm:"size:255;not null;"`
	Name         string `gorm:"size:255;not null;"`
	AvatarURL    string `gorm:"size:5000"`
	RefreshToken string `gorm:"size:255;"`
	//
	// gorm model
	//
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
}
