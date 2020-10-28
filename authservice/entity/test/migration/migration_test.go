package migration

import (
	"authservice/entity"
	"authservice/pckg/storage"
	"testing"
)

var (
	postgres = storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		"postgres",
		"toster123",
		"pet_finder_user",
		"5432",
		"disable",
	)
)

func TestGormMigrationFlow(t *testing.T) {
	err := entity.GORMMigration(postgres.DB)
	if err != nil {
		t.Fatal(err)
	}
}
