package app

import (
	"dialogservice/pckg/conf"
	"dialogservice/pckg/storage"
	"dialogservice/repository"
	"dialogservice/service"
	"github.com/spf13/viper"
)

var (
	testRoot                = "C:/PetFinderRepos/petfinder/dialogservice"
	testConfigs             map[string]*viper.Viper
	testPostgresOrm         *storage.Storage
	testRepositoryDialogAPI repository.DialogRepositoryAPI
	testServiceDialogAPI    *service.DialogServiceAPI
	testApplication         *Application
)

func initFields() {
	testRoot = "C:/PetFinderRepos/petfinder/authservice"
	testConfigs = conf.ReadConfigs(
		testRoot,
		"app",
		"event_receivers",
	)
	testPostgresOrm = storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		"postgres",
		"toster123",
		"pet_finder_user",
		"5432",
		"disable",
	)
	testRepositoryDialogAPI = repository.NewGormDialogRepositoryAPI(testPostgresOrm.DB)
	testServiceDialogAPI = service.NewDialogServiceAPI()
	testApplication = newTestApp(
		testConfigs,
		testServiceDialogAPI,
		testRepositoryDialogAPI,
	)
}