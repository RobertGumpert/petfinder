package entity

import "gorm.io/gorm"

func GORMMigration(db *gorm.DB) error {
	var(
		user User
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.AutoMigrate(user); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
