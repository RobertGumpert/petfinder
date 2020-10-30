package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fileservice/pckg/runtimeinfo"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func serverStart() {
	engine := gin.Default()
	uploadRouter := engine.Group("/upload")
	{
		avatarRouter := uploadRouter.Group("/avatar")
		{
			avatarRouter.POST("", func(context *gin.Context) {
				resp, err := uploadFileFormData(context, avatarGroupFiles)
				if err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.AbortWithStatusJSON(http.StatusOK, resp)
				return
			})
		}
		advertRouter := uploadRouter.Group("/advert")
		{
			advertRouter.POST("", func(context *gin.Context) {
				resp, err := uploadFileFormData(context, advertGroupFiles)
				if err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.AbortWithStatusJSON(http.StatusOK, resp)
				return
			})
		}
	}
	downloadRouter := engine.Group("/download")
	{
		avatarRouter := downloadRouter.Group("/avatar")
		{
			avatarRouter.POST("/base64", func(context *gin.Context) {
				response, err := downloadFile64(context, avatarGroupFiles)
				if err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.AbortWithStatusJSON(http.StatusOK, response)
			})
			avatarRouter.GET("/id/:id", func(context *gin.Context) {
				if err := download(context, avatarGroupFiles); err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.AbortWithStatus(http.StatusOK)
				return
			})
		}
		advertRouter := downloadRouter.Group("/advert")
		{
			advertRouter.POST("/base64", func(context *gin.Context) {
				response, err := downloadFile64(context, advertGroupFiles)
				if err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.AbortWithStatusJSON(http.StatusOK, response)
			})
			advertRouter.GET("/id/:id", func(context *gin.Context) {
				if err := download(context, advertGroupFiles); err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.AbortWithStatus(http.StatusOK)
				return
			})
		}
	}
	port := configs["app"].GetString("port")
	if !strings.Contains(port, ":") {
		port = strings.Join([]string{
			":",
			port,
		}, "")
	}
	err := engine.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

func download(context *gin.Context, fileGroup groupFiles) error {
	paramId := context.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return err
	}
	file, err := findFileByID(fileGroup, uint64(id))
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	buffer := make([]byte, info.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	fileContentType := http.DetectContentType(buffer)
	context.Writer.Header().Set("Content-Disposition", "attachment; filename="+info.Name())
	context.Writer.Header().Set("Content-Type", fileContentType)
	context.Writer.Header().Set("Content-Length", strconv.Itoa(int(info.Size())))
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = io.Copy(context.Writer, file)
	return err
}

func downloadFile64(context *gin.Context, fileGroup groupFiles) (*ListFile64ViewModel, error) {
	viewModel := new(DownloadViewModel)
	jsonForm := context.PostForm("json")
	err := json.Unmarshal([]byte(jsonForm), viewModel)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	files := readFilesByID(fileGroup, viewModel)
	response := filesToBase64(files)
	return response, nil
}

func uploadFileFormData(context *gin.Context, fileGroup groupFiles) (*URLViewModel, error) {
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, 2<<20)
	viewModel := new(UploadViewModel)
	jsonForm := context.PostForm("json")
	err := json.Unmarshal([]byte(jsonForm), viewModel)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	header, err := context.FormFile("file")
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	file, err := header.Open()
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	defer file.Close()
	buffer := make([]byte, header.Size)
	_, err = file.Read(buffer)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	contentType := http.DetectContentType(buffer)
	var fileType typeFiles = 0
	if ft, exist := imageTypes[contentType]; exist {
		fileType = ft
	} else {
		err = errors.New("Non valid file type. ")
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	reader := bufio.NewReader(bytes.NewBuffer(buffer))
	if err := saveFile(reader, fileGroup, fileType, viewModel); err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	group := ""
	switch fileGroup {
	case avatarGroupFiles:
		group = "avatar/id/"
		break
	case advertGroupFiles:
		group = "advert/id/"
		break
	}
	url := strings.Join([]string{
		configs["app"].GetString("addr"),
		":",
		configs["app"].GetString("port"),
		"/download/",
		group,
		strconv.Itoa(int(viewModel.ID)),
	}, "")
	return &URLViewModel{URL: url}, err
}
