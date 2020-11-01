package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

type apiHttpHandler struct{}

func newApiHttpHandler() *apiHttpHandler {
	return &apiHttpHandler{}
}


func (a *apiHttpHandler) getServer() func() {
	port := application.configs["app"].GetString("port")
	if !strings.Contains(port, ":") {
		port = strings.Join([]string{":", port}, "")
	}
	engine := gin.Default()

	return func() {
		err := engine.Run(port)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

