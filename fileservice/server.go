package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fileservice/pckg/runtimeinfo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func serverStart() {
	engine := gin.Default()
	uploadRouter := engine.Group("/upload")
	{
		avatarRouter := uploadRouter.Group("/avatar")
		{
			avatarRouter.POST("", func(context *gin.Context) {
				err := uploadFileFormData(context, avatarGroupFiles)
				if err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				return
			})
		}
		advertRouter := uploadRouter.Group("/advert")
		{
			advertRouter.POST("", func(context *gin.Context) {
				err := uploadFileFormData(context, advertGroupFiles)
				if err != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
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

func downloadFile64(context *gin.Context, fileGroup groupFiles) (*File64ListViewModel, error) {
	viewModel := new(DownloadViewModel)
	jsonForm := context.PostForm("json")
	err := json.Unmarshal([]byte(jsonForm), viewModel)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return nil, err
	}
	files := readFiles(fileGroup, viewModel)
	response := filesToBase64(files)
	return response, nil
}

func uploadFileFormData(context *gin.Context, fileGroup groupFiles) error {
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, 2<<20)
	viewModel := new(UploadViewModel)
	jsonForm := context.PostForm("json")
	err := json.Unmarshal([]byte(jsonForm), viewModel)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return err
	}
	header, err := context.FormFile("file")
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return err
	}
	file, err := header.Open()
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return err
	}
	defer file.Close()
	buffer := make([]byte, header.Size)
	_, err = file.Read(buffer)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return err
	}
	contentType := http.DetectContentType(buffer)
	var fileType typeFiles = 0
	if ft, exist := imageTypes[contentType]; exist {
		fileType = ft
	} else {
		err = errors.New("Non valid file type. ")
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return err
	}
	reader := bufio.NewReader(bytes.NewBuffer(buffer))
	if err := saveFile(reader, fileGroup, fileType, viewModel); err != nil {
		log.Println(runtimeinfo.Runtime(1), "ERROR=[", err, "]")
		return err
	}
	return nil
}
