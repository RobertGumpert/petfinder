package app

import (
	"dialogservice/mapper"
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

func (a *apiHttpHandler) getServer() (*gin.Engine, func()) {
	port := application.configs["app"].GetString("port")
	if !strings.Contains(port, ":") {
		port = strings.Join([]string{":", port}, "")
	}
	engine := gin.Default()
	api := engine.Group("/api/user", a.middlewareAccessToken)
	{
		dialogs := api.Group("/dialog")
		{
			dialogs.GET("/get", a.downloadDialogs)
			dialogs.POST("/create", a.createNewDialog)
		}
		messages := api.Group("/message")
		{
			messages.POST("/send", a.addNewMessage)
			messages.POST("/batching/next", a.downloadNextMessagesBatch)
		}
	}
	return engine, func() {
		err := engine.Run(port)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func (a *apiHttpHandler) getMiddlewareToken(ctx *gin.Context, header, key string) (token string, err error) {
	token = ""
	err = errors.New("token is nil")
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
	authUser, err := applicationHttpRequests.isAuthorized(token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("authorization", authUser)
}

func (a *apiHttpHandler) createNewDialog(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	receiverUser := new(mapper.UserViewModel)
	if err := ctx.BindJSON(receiverUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	newDialogViewModel, err := application.dialogServiceAPI.CreateNewDialog(
		authUser,
		receiverUser,
		application.dialogAPIPostgresRepository,
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
	ctx.AbortWithStatusJSON(http.StatusOK, newDialogViewModel)
	return
}

func (a *apiHttpHandler) downloadDialogs(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	downloadDialogsViewModel, err := application.dialogServiceAPI.DownloadDialogs(
		authUser,
		application.dialogAPIPostgresRepository,
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
	ctx.AbortWithStatusJSON(http.StatusOK, downloadDialogsViewModel)
	return
}

func (a *apiHttpHandler) downloadNextMessagesBatch(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	batchingViewModel := new(mapper.NextMessagesBatchViewModel)
	if err := ctx.BindJSON(batchingViewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	batchingViewModel.UserReceiver = authUser
	response, err := application.dialogServiceAPI.DownloadNextMessagesBatch(
		batchingViewModel,
		application.dialogAPIPostgresRepository,
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
	ctx.AbortWithStatusJSON(http.StatusOK, response)
	return
}

func (a *apiHttpHandler) addNewMessage(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	addMessageViewModel := new(mapper.AddNewMessageViewModel)
	if err := ctx.BindJSON(addMessageViewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	addMessageViewModel.UserReceiver = authUser
	response, err := application.dialogServiceAPI.AddNewMessage(
		addMessageViewModel,
		application.dialogAPIPostgresRepository,
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
	ctx.AbortWithStatusJSON(http.StatusOK, response)
	return
}
