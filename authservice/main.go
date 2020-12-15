package main

import (
	"authservice/app"
	"authservice/pckg/conf"
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
		"event_receivers",
	)
	application := app.NewApp(configs)
	application.HttpServerRun()
}
