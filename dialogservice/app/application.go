package app

import (
	"dialogservice/repository"
	"dialogservice/service"
	"github.com/gin-gonic/gin"
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
	HttpAPI                     *gin.Engine
	HttpServerRun               func()
}

func NewApp(configs map[string]*viper.Viper) *Application {
	application = new(Application)
	application.configs = configs
	postgres := postgresInit(true)
	application.dialogAPIPostgresRepository = repository.NewGormDialogRepositoryAPI(postgres)
	application.dialogServiceAPI = service.NewDialogServiceAPI()
	//
	applicationHttpApi = newApiHttpHandler()
	ginEngine, httpServerRun := applicationHttpApi.getServer()
	application.HttpAPI = ginEngine
	application.HttpServerRun = httpServerRun
	applicationHttpRequests = newHttpRequests()
	//
	return application
}

func newTestApp(configs map[string]*viper.Viper, dsApi *service.DialogServiceAPI, dbApi repository.DialogRepositoryAPI) *Application {
	application = new(Application)
	application.configs = configs
	//
	application.dialogAPIPostgresRepository = dbApi
	application.dialogServiceAPI = dsApi
	//
	applicationHttpApi = newApiHttpHandler()
	ginEngine, httpServerRun := applicationHttpApi.getServer()
	application.HttpAPI = ginEngine
	application.HttpServerRun = httpServerRun
	applicationHttpRequests = newHttpRequests()
	//
	return application
}
