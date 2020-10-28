package app

import (
	"authservice/service"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	api := engine.Group("/api/user")
	{
		api.POST("/register", a.register)
		api.POST("/authorized", a.authorized)

		auth := api.Group("/access", a.middlewareAccessToken)
		{
			auth.POST("", a.isAuthorized)
			auth.POST("/update", a.updateAccessToken)
		}

		passReset := api.Group("/password")
		{
			passReset.POST("/token", a.getResetPasswordToken)
			passReset.POST("/reset", a.resetPassword)
		}

		update := api.Group("/update", a.middlewareAccessToken)
		{
			update.POST("", a.update)
		}

		getting := api.Group("/get")
		{
			getting.POST("", a.get)
		}
	}
	return func() {
		err := engine.Run(port)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func (a *apiHttpHandler) getMiddlewareToken(ctx *gin.Context, header, key string) (token string, err error) {
	//
	token = ""
	err = errors.New("token is nil")
	//
	headerValue := ctx.GetHeader(header)
	if headerValue == "" {
		return token, err
	}
	headerValueSplit := strings.Split(headerValue, " ")
	if len(headerValueSplit) != 2 {
		return token, err
	}
	if headerValueSplit[0] != key {
		return token, err
	}
	token = headerValueSplit[1]
	return token, nil
}

func (a *apiHttpHandler) middlewareAccessToken(ctx *gin.Context) {
	token, err := a.getMiddlewareToken(
		ctx,
		"Authorization",
		"Bearer",
	)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("authorization", token)
}

func (a *apiHttpHandler) register(ctx *gin.Context) {
	json := new(service.RegisterUserViewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	response, err := application.userService.Register(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	go applicationHttpRequests.mailerAuthorized(
		response.Name,
		response.Telephone,
		response.Email,
		true,
	)
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (a *apiHttpHandler) authorized(ctx *gin.Context) {
	json := new(service.AuthorizationUserViewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	access, _, response, err := application.userService.Authorized(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	go applicationHttpRequests.mailerAuthorized(
		response.Name,
		response.Telephone,
		response.Email,
		false,
	)
	ctx.AbortWithStatusJSON(http.StatusOK, &struct {
		Token string                 `json:"token"`
		User  *service.UserViewModel `json:"user"`
	}{
		Token: access,
		User:  response,
	})
}

func (a *apiHttpHandler) isAuthorized(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	json := new(service.IsAuthorizedViewModel)
	json.Access = token
	response, err := application.userService.IsAuthorized(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (a *apiHttpHandler) updateAccessToken(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	json := new(service.NewAccessTokenViewModel)
	json.Access = token
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	access, response, err := application.userService.UpdateAccessToken(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, &struct {
		Token string                 `json:"token"`
		User  *service.UserViewModel `json:"user"`
	}{
		Token: access,
		User:  response,
	})
}

func (a *apiHttpHandler) getResetPasswordToken(ctx *gin.Context) {
	json := new(service.ResetPasswordViewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	token, err := application.userService.GetResetPasswordToken(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	go applicationHttpRequests.mailerResetPasswordToken(token, json.Email)
	ctx.AbortWithStatus(http.StatusOK)
}

func (a *apiHttpHandler) resetPassword(ctx *gin.Context) {
	json := new(service.ResetPasswordViewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	access, _, err := application.userService.ResetPassword(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	go applicationHttpRequests.mailerResetPassword(json.Email)
	ctx.AbortWithStatusJSON(http.StatusOK, &struct {
		Token string `json:"token"`
	}{
		Token: access,
	})
}

func (a *apiHttpHandler) update(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	user, err := application.userService.IsAuthorized(
		&service.IsAuthorizedViewModel{Access: token},
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	json := new(service.UpdateUserViewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	response, err := application.userService.Update(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	if user.Name != response.Name {
		applicationHttpRequests.eventUserUpdateName(response.UserID, response.Name)
	}
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (a *apiHttpHandler) get(ctx *gin.Context) {
	json := new(service.FindUserViewModel)
	if err := ctx.BindJSON(json); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: service.ErrorNonValidData.Error(),
		})
		return
	}
	response, err := application.userService.Get(
		json,
		application.userPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}