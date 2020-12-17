package app

import (
	"advertservice/mapper"
	"advertservice/pckg/runtimeinfo"
	"advertservice/repository"
	"advertservice/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"strconv"
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
	HttpAPI                   *gin.Engine
	HttpServerRun             func()
}

func NewApp(configs map[string]*viper.Viper) *Application {
	application = new(Application)
	application.configs = configs
	postgres := postgresInit(true)
	application.advertPostgresRepository = repository.NewGormAdvertRepository(postgres)
	radius, err := strconv.ParseFloat(configs["app"].GetString("radius"), 64)
	if err != nil {
		radius = float64(1)
		go log.Println(runtimeinfo.Runtime(1), " ERROR: [", err, "]")
	}
	application.advertPostgresSearchModel = repository.NewGormSquareSearchModel(
		postgres,
		mapper.CompareAdvertTime,
		mapper.OneKilometerScale*radius,
	)
	application.advertService = service.NewAdvertService(
		mapper.LifetimeOfFoundAnimalAdvert,
		mapper.LifetimeOfLostAnimalAdvert,
		mapper.CompareAdvertTime,
	)
	//
	applicationHttpApi = newApiHttpHandler()
	ginEngine, httpServerRun := applicationHttpApi.getServer()
	application.HttpAPI = ginEngine
	application.HttpServerRun = httpServerRun
	applicationHttpRequests = newHttpRequests()
	//
	return application
}

func newTestApp(configs map[string]*viper.Viper, as *service.AdvertService, db repository.AdvertRepository, sm repository.SearchModel) *Application {
	application = new(Application)
	application.configs = configs
	//
	application.advertPostgresRepository = db
	application.advertPostgresSearchModel = sm
	application.advertService = as
	//
	applicationHttpApi = newApiHttpHandler()
	ginEngine, httpServerRun := applicationHttpApi.getServer()
	application.HttpAPI = ginEngine
	application.HttpServerRun = httpServerRun
	applicationHttpRequests = newHttpRequests()
	//
	return application
}
