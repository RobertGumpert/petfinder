package main

import (
	"authservice/app"
	"authservice/pckg/runtimeinfo"
	"github.com/spf13/viper"
	"log"
	"path"
	"runtime"
)

func readConfigs(files ...string) map[string]*viper.Viper {
	_, file, _, _ := runtime.Caller(0)
	root := path.Dir(file)
	configs := make(map[string]*viper.Viper)
	var read = func(name string) *viper.Viper {
		vpr := viper.New()
		vpr.SetConfigFile(root + "/" + name + ".yaml")
		if err := vpr.ReadInConfig(); err != nil {
			log.Fatal(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		}
		return vpr
	}
	for _, file := range files {
		configs[file] = read(file)
	}
	return configs
}

func main() {
	configs := readConfigs(
		"app",
		"event_receivers",
	)
	app.NewApp(configs)
}
