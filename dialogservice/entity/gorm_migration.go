package entity

import "gorm.io/gorm"

func GORMMigration(db *gorm.DB) error {
	var (
		dialog     DialogEntity
		dialogUser DialogUserEntity
		message    MessageEntity
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.AutoMigrate(dialog, dialogUser, message); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
