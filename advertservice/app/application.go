package app

import (
	"advertservice/mapper"
	"advertservice/repository"
	"advertservice/service"
	"github.com/spf13/viper"
)

var (
	application             *Application
	applicationHttpApi      *apiHttpHandler
	applicationHttpRequests *httpRequests
)

type Application struct {
	configs                   map[string]*viper.Viper
	advertPostgresRepository  repository.AdvertRepository
	advertPostgresSearchModel repository.SearchModel
	advertService             *service.AdvertService
}

func NewApp(configs map[string]*viper.Viper) *Application {
	application = new(Application)
	application.configs = configs
	postgres := postgresInit(true)
	application.advertPostgresRepository = repository.NewGormAdvertRepository(postgres)
	application.advertPostgresSearchModel = repository.NewGormSquareSearchModel(postgres, mapper.CompareAdvertTime, mapper.OneKilometerScale)
	application.advertService = service.NewAdvertService(
		mapper.LifetimeOfFoundAnimalAdvert,
		mapper.LifetimeOfLostAnimalAdvert,
		mapper.CompareAdvertTime,
	)
	//
	applicationHttpApi = newApiHttpHandler()
	httpServerRun := applicationHttpApi.getServer()
	//
	applicationHttpRequests = newHttpRequests()
	//
	httpServerRun()
	return application
}
