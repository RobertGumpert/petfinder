package main

import (
	"github.com/spf13/viper"
	"log"
	"path"
	"runtime"
)

type groupFiles uint64
type typeFiles uint64

const (
	avatarGroupFiles groupFiles = 0
	advertGroupFiles groupFiles = 1
	jpegFileType     typeFiles  = 10
	pngFileType      typeFiles  = 11
)

var (
	configs    map[string]*viper.Viper
	root       string
	imageTypes = map[string]typeFiles{
		"image/jpeg": jpegFileType,
		"image/png":  pngFileType,
	}
)

func main() {
	configs = readConfigs(
		"app",
	)
	serverStart()
}

func readConfigs(files ...string) map[string]*viper.Viper {
	setRoot()
	configs := make(map[string]*viper.Viper)
	var read = func(name string) *viper.Viper {
		vpr := viper.New()
		vpr.SetConfigFile(root + "/" + name + ".yaml")
		if err := vpr.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
		return vpr
	}
	for _, file := range files {
		configs[file] = read(file)
	}
	return configs
}

func setRoot() {
	_, file, _, _ := runtime.Caller(0)
	root = path.Dir(file)
}