package entity

import "gorm.io/gorm"

func GORMMigration(db *gorm.DB) error {
	var(
		advert Advert
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.AutoMigrate(advert); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}