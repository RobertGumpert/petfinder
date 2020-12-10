package app

import (
	"dialogservice/repository"
	"github.com/spf13/viper"
)

var (
	application *Application
)

type Application struct {
	configs                     map[string]*viper.Viper
	dialogAPIPostgresRepository repository.DialogRepositoryAPI
}

func NewApplication(configs map[string]*viper.Viper) *Application {
	return &Application{configs: configs}
}
