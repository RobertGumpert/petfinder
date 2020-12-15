package conf

import (
	"dialogservice/pckg/runtimeinfo"
	"github.com/spf13/viper"
	"log"
)

func ReadConfigs(root string, files ...string) map[string]*viper.Viper {
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
