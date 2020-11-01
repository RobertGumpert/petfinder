package app

import (
	"advertservice/repository"
	"github.com/spf13/viper"
)

var (
	application             *Application
	applicationHttpApi      *apiHttpHandler
	applicationHttpRequests *httpRequests
)

type Application struct {
	configs                  map[string]*viper.Viper
	advertPostgresRepository repository.AdvertRepository
}

func NewApp(configs map[string]*viper.Viper) *Application {
	application = new(Application)
	application.configs = configs
	application.advertPostgresRepository = repository.NewAdvertGormRepository(postgresInit(true))
	//
	applicationHttpApi = newApiHttpHandler()
	httpServerRun := applicationHttpApi.getServer()
	//
	applicationHttpRequests = newHttpRequests()
	//
	httpServerRun()
	return application
}
