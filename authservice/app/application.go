package app

import (
	"authservice/repository"
	"authservice/service"
	"github.com/spf13/viper"
	"time"
)

var (
	application             *Application
	applicationHttpApi      *apiHttpHandler
	applicationHttpRequests *httpRequests
)

type Application struct {
	configs                map[string]*viper.Viper
	userPostgresRepository repository.UserRepository
	userService            *service.User
}

func NewApp(configs map[string]*viper.Viper) *Application {
	application = new(Application)
	application.configs = configs
	application.userPostgresRepository = repository.NewUserGormRepository(postgresInit(true))
	application.userService = service.NewUserService(
		[]byte("auth_service"),
		30*time.Minute,
		365*7*24*time.Hour,
		30*time.Minute,
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
