package app

import (
	"advertservice/pckg/storage"
	"gorm.io/gorm"
)

func postgresInit(migrate bool) *gorm.DB {
	postgres := storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		application.configs["app"].GetString("postgres_username"),
		application.configs["app"].GetString("postgres_password"),
		application.configs["app"].GetString("postgres_name"),
		application.configs["app"].GetString("postgres_port"),
		application.configs["app"].GetString("postgres_ssl"),
	)
	if migrate {
		if err := entity.GORMMigration(postgres.DB); err != nil {
			log.Fatal(err)
		}
	}
	return postgres.DB
}