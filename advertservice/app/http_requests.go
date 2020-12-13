package app

import (
	"advertservice/mapper"
	"advertservice/pckg/runtimeinfo"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type httpRequests struct{}

func newHttpRequests() *httpRequests {
	return &httpRequests{}
}

func (h *httpRequests) isAuthorized(token string) (*mapper.UserViewModel, error) {
	url := strings.Join([]string{
		"http://",
		application.configs["app"].GetString("auth_service"),
		"/api/user/access",
	}, "")
	client := http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	req.Header.Set("Authorization", strings.Join([]string{"Bearer", token}, " "))
	res, err := client.Do(req)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Http code : " + res.Status)
	}
	authUser := new(mapper.UserViewModel)
	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	if err := json.Unmarshal(bts, authUser); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	return authUser, nil
}

func (h *httpRequests) saveImage(id uint64, file multipart.File, fileHeader *multipart.FileHeader) (*http.Response, string, error) {
	jsonBuffer, err := json.Marshal(&struct {
		ID                       uint64 `json:"id"`
		AdditionalIdentification string `json:"additional_identification"`
	}{
		ID: id,
	})
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	var requestBuffer bytes.Buffer
	fileBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	requestWriter := multipart.NewWriter(&requestBuffer)
	var fileWriter io.Writer
	var fileReader io.Reader = bufio.NewReader(bytes.NewBuffer(fileBuffer))
	mimeHeader := make(textproto.MIMEHeader)
	mimeHeader.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"file", fileHeader.Filename))
	mimeHeader.Set("Content-Type", "application/octet-stream")
	if fileWriter, err = requestWriter.CreatePart(mimeHeader); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	if _, err = io.Copy(fileWriter, fileReader); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	var jsonWriter io.Writer
	var jsonReader io.Reader = bufio.NewReader(bytes.NewBuffer(jsonBuffer))
	if jsonWriter, err = requestWriter.CreateFormField("json"); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	if _, err = io.Copy(jsonWriter, jsonReader); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	err = requestWriter.Close()
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	ur := strings.Join([]string{
		"http://",
		application.configs["app"].GetString("file_service"),
		"/upload/advert",
	}, "")
	client := http.Client{}
	req, err := http.NewRequest("POST", ur, &requestBuffer)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	req.Header.Set("Content-Type", requestWriter.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, "", errors.New("Bad request. ")
	}
	var jsonResponse = struct {
		URL string `json:"url"`
	}{}
	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	if err := json.Unmarshal(bts, &jsonResponse); err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, "", err
	}
	return res, jsonResponse.URL, nil
}
