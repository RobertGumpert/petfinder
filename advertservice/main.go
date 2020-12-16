package main

import (
	"advertservice/app"
	"advertservice/pckg/conf"
	"path"
	"runtime"
)

var root string

func main() {
	_, file, _, _ := runtime.Caller(0)
	root = path.Dir(file)
	configs := conf.ReadConfigs(
		root,
		"app",
	)
	application := app.NewApp(configs)
	application.HttpServerRun()
}
