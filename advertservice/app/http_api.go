package app

import (
	"advertservice/mapper"
	"advertservice/pckg/runtimeinfo"
	"encoding/json"
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
	api := engine.Group("/api/advert")
	{
		middleware := api.Group("/user", a.middlewareAccessToken)
		{
			middleware.POST("/add", a.addAdvert)
			middleware.POST("/update", a.updateAdvert)
			middleware.POST("/list", a.userAdverts)
			middleware.POST("/close", a.closeAdvert)
			middleware.POST("/refresh", a.refreshAdvert)
		}
		getting := api.Group("/get")
		{
			getting.POST("/in/area", a.searchInArea)
		}
	}
	event := engine.Group("/event")
	{
		userEvents := event.Group("/user")
		{
			userEvents.POST("/update/name", a.updateUser)
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
	var authUser *mapper.UserViewModel
	if authUser, err = applicationHttpRequests.isAuthorized(token); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("authorization", authUser)
}

func (a *apiHttpHandler) addAdvert(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	if err := authUser.Validator(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	//
	viewModel := new(mapper.CreateAdvertViewModel)
	//
	jsonForm := ctx.PostForm("json")
	if jsonForm != "" {
		err := json.Unmarshal([]byte(jsonForm), viewModel)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: mapper.ErrorNonValidData.Error(),
			})
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	viewModel.AdOwnerID = authUser.UserID
	viewModel.AdOwnerName = authUser.Name
	response, err := application.advertService.CreateAdvert(viewModel, application.advertPostgresRepository, nil)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	_, downloadUrl, err := applicationHttpRequests.saveImage(response.AdID, file, fileHeader)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := application.advertService.UpdateImage(
		&mapper.UpdateImageViewModel{AdID: response.AdID, ImageUrl: downloadUrl},
		application.advertPostgresRepository,
		nil,
	); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	response.ImageUrl = downloadUrl
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (a *apiHttpHandler) updateAdvert(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	if err := authUser.Validator(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	//
	viewModel := new(mapper.UpdateAdvertViewModel)
	jsonForm := ctx.PostForm("json")
	if jsonForm != "" {
		err := json.Unmarshal([]byte(jsonForm), viewModel)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: mapper.ErrorNonValidData.Error(),
			})
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	update, err := application.advertService.Update(
		viewModel,
		application.advertPostgresRepository,
		nil,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	downloadUrl := ""
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err == nil {
		_, downloadUrl, err = applicationHttpRequests.saveImage(viewModel.AdID, file, fileHeader)
		if err != nil {
			log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	if strings.TrimSpace(downloadUrl) != "" {
		if err := application.advertService.UpdateImage(
			&mapper.UpdateImageViewModel{AdID: update.AdID, ImageUrl: downloadUrl},
			application.advertPostgresRepository,
			nil,
		); err != nil {
			log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})
			return
		}
		update.ImageUrl = downloadUrl
	}
	ctx.AbortWithStatusJSON(http.StatusOK, update)
}

func (a *apiHttpHandler) userAdverts(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	if err := authUser.Validator(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	//
	viewModel := new(mapper.IdentifierOwnerViewModel)
	viewModel.AdOwnerName = authUser.Name
	viewModel.AdOwnerID = authUser.UserID
	//
	list, err := application.advertService.ListMyAdverts(
		viewModel,
		application.advertPostgresRepository,
		nil,
	)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, list)
}

func (a *apiHttpHandler) refreshAdvert(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	if err := authUser.Validator(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	//
	viewModel := new(mapper.IdentifierAdvertViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	update, err := application.advertService.RefreshAdvert(
		viewModel,
		application.advertPostgresRepository,
		nil,
	)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, update)
}

func (a *apiHttpHandler) closeAdvert(ctx *gin.Context) {
	authUser := ctx.MustGet("authorization").(*mapper.UserViewModel)
	if err := authUser.Validator(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	//
	viewModel := new(mapper.IdentifierAdvertViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	update, err := application.advertService.CloseAdvert(
		viewModel,
		application.advertPostgresRepository,
		nil,
	)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, update)
}

func (a *apiHttpHandler) searchInArea(ctx *gin.Context) {
	viewModel := new(mapper.SearchInAreaViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	list, err := application.advertService.SearchInArea(
		viewModel,
		application.advertPostgresSearchModel,
		nil,
	)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, list)
}

func (a *apiHttpHandler) updateUser(ctx *gin.Context) {
	originJson := new(struct {
		UserID uint64 `json:"user_id"`
		Name   string `json:"name"`
	})
	if err := ctx.BindJSON(originJson); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	if err := application.advertService.UpdateOwnerName(
		&mapper.IdentifierOwnerViewModel{
			AdOwnerID:   originJson.UserID,
			AdOwnerName: originJson.Name,
		},
		application.advertPostgresRepository,
		nil,
	); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
