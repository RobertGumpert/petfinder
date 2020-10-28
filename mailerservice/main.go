package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"mailerservice/mailer"
	"net/http"
	"strings"
)

var (
	postman *mailer.Postman
	configs map[string]*viper.Viper
)

func main() {
	configs = readConfigs(
		"app",
		"mail",
	)
	port := configs["app"].GetString("port")
	if !strings.Contains(port, ":") {
		port = strings.Join([]string{":", port}, "")
	}
	engine := gin.Default()
	engine.POST("/user/register", userRegister)
	engine.POST("/user/authorized", userAuthorized)
	engine.POST("/user/pass/token", resetPasswordToken)
	engine.POST("/user/pass/reset", resetPassword)
	err := engine.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

func userRegister(ctx *gin.Context) {
	type viewModel struct {
		Telephone string `json:"telephone"`
		Email     string `json:"email"`
		Name      string `json:"name"`
	}
	json := new(viewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	err := postman.SendLetter(
		"auth_service",
		"sign_up",
		json.Email,
		json,
		mailer.TCP,
	)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func userAuthorized(ctx *gin.Context) {
	type viewModel struct {
		Telephone string `json:"telephone"`
		Email     string `json:"email"`
		Name      string `json:"name"`
	}
	json := new(viewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	err := postman.SendLetter(
		"auth_service",
		"sign_in",
		json.Email,
		json,
		mailer.TCP,
	)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func resetPasswordToken(ctx *gin.Context) {
	type viewModel struct {
		Token string `json:"token"`
		Email string `json:"email"`
	}
	json := new(viewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	err := postman.SendLetter(
		"auth_service",
		"reset_pass",
		json.Email,
		json,
		mailer.TCP,
	)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func resetPassword(ctx *gin.Context) {
	type viewModel struct {
		Token string `json:"token"`
		Email string `json:"email"`
	}
	json := new(viewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	err := postman.SendLetter(
		"auth_service",
		"reset_pass_2",
		json.Email,
		json,
		mailer.TCP,
	)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
