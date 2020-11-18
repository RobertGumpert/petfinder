package app

import (
	"authservice/mapper"
	"authservice/pckg/runtimeinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
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
			update.POST("/avatar", a.updateAvatar)
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
	viewModel := new(mapper.RegisterUserViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	response, err := application.userService.Register(
		viewModel,
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
	viewModel := new(mapper.AuthorizationUserViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	access, _, response, err := application.userService.Authorized(
		viewModel,
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
		Token string                `json:"token"`
		User  *mapper.UserViewModel `json:"user"`
	}{
		Token: access,
		User:  response,
	})
}

func (a *apiHttpHandler) isAuthorized(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	viewModel := new(mapper.IsAuthorizedViewModel)
	viewModel.Access = token
	response, err := application.userService.IsAuthorized(
		viewModel,
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
	viewModel := new(mapper.NewAccessTokenViewModel)
	viewModel.Access = token
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	access, response, err := application.userService.UpdateAccessToken(
		viewModel,
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
		Token string                `json:"token"`
		User  *mapper.UserViewModel `json:"user"`
	}{
		Token: access,
		User:  response,
	})
}

func (a *apiHttpHandler) getResetPasswordToken(ctx *gin.Context) {
	viewModel := new(mapper.ResetPasswordViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	token, err := application.userService.GetResetPasswordToken(
		viewModel,
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
	go applicationHttpRequests.mailerResetPasswordToken(token, viewModel.Email)
	ctx.AbortWithStatus(http.StatusOK)
}

func (a *apiHttpHandler) resetPassword(ctx *gin.Context) {
	viewModel := new(mapper.ResetPasswordViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	access, _, err := application.userService.ResetPassword(
		viewModel,
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
	go applicationHttpRequests.mailerResetPassword(viewModel.Email)
	ctx.AbortWithStatusJSON(http.StatusOK, &struct {
		Token string `json:"token"`
	}{
		Token: access,
	})
}

func (a *apiHttpHandler) update(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	user, err := application.userService.IsAuthorized(
		&mapper.IsAuthorizedViewModel{Access: token},
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
	viewModel := new(mapper.UpdateUserViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	response, err := application.userService.Update(
		viewModel,
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

func (a *apiHttpHandler) updateAvatarProxy(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	user, err := application.userService.IsAuthorized(
		&mapper.IsAuthorizedViewModel{Access: token},
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
	u := &url.URL{
		Scheme: "http",
		Host:   application.configs["app"].GetString("file_service"),
		Path:   "/upload/avatar/id/" + strconv.Itoa(int(user.UserID)),
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.PostForm = ctx.Request.PostForm
		req.Host = u.Host
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = u.Path
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (a *apiHttpHandler) updateAvatar(ctx *gin.Context) {
	token := ctx.MustGet("authorization").(string)
	user, err := application.userService.IsAuthorized(
		&mapper.IsAuthorizedViewModel{Access: token},
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
	func(context *gin.Context, id uint64, app *Application, req *httpRequests) {
		_, imageUrl, err := req.saveAvatar(context, id)
		if err == nil {
			err := app.userService.UpdateAvatar(&mapper.UpdateAvatarViewModel{
				UserID:    id,
				AvatarUrl: imageUrl,
			}, app.userPostgresRepository, nil)
			if err != nil {
				log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
				return
			}
		} else {
			log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
			return
		}
	}(ctx, user.UserID, application, applicationHttpRequests)
	ctx.AbortWithStatus(http.StatusOK)
}

func (a *apiHttpHandler) get(ctx *gin.Context) {
	viewModel := new(mapper.FindUserViewModel)
	if err := ctx.BindJSON(viewModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{
			Error: mapper.ErrorNonValidData.Error(),
		})
		return
	}
	response, err := application.userService.Get(
		viewModel,
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
