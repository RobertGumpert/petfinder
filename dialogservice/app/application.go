package app

import (
	"dialogservice/repository"
	"dialogservice/service"
	"github.com/spf13/viper"
)

var (
	application             *Application
	applicationHttpApi      *apiHttpHandler
	applicationHttpRequests *httpRequests
)

type Application struct {
	configs                     map[string]*viper.Viper
	dialogAPIPostgresRepository repository.DialogRepositoryAPI
	dialogServiceAPI            *service.DialogServiceAPI
}

func NewApp(configs map[string]*viper.Viper) *Application {
	application = new(Application)
	application.configs = configs
	postgres := postgresInit(true)
	application.dialogAPIPostgresRepository = repository.NewGormDialogRepositoryAPI(postgres)
	application.dialogServiceAPI = service.NewDialogServiceAPI()
	//
	applicationHttpApi = newApiHttpHandler()
	httpServerRun := applicationHttpApi.getServer()
	//
	applicationHttpRequests = newHttpRequests()
	//
	httpServerRun()
	return application
}
