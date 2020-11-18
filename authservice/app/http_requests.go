package app

import (
	"authservice/pckg/runtimeinfo"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

func (h *httpRequests) eventUserUpdateName(userId uint64, userName string) {
	js, err := json.Marshal(&struct {
		UserID uint64 `json:"user_id"`
		Name   string `json:"name"`
	}{
		UserID: userId,
		Name:   userName,
	})
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
		return
	}
	for _, value := range application.configs["event_receivers"].AllSettings() {
		serviceName := value.(string)
		go func(addr string, body []byte) {
			url := strings.Join([]string{
				"http://",
				addr,
				"/event/user/update/name",
			}, "")
			response, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil || response.StatusCode != http.StatusOK {
				log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
			}
		}(serviceName, js)
	}
}


func (h *httpRequests) saveAvatar(context *gin.Context, id uint64) (*http.Response, string, error) {
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
	file, fileHeader, err := context.Request.FormFile("file")
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
		"/upload/avatar",
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

func (h *httpRequests) mailerAuthorized(name, telephone, email string, isRegister bool) {
	endpoint := "register"
	if !isRegister {
		endpoint = "authorized"
	}
	body, err := json.Marshal(&struct {
		Telephone string `json:"telephone"`
		Email     string `json:"email"`
		Name      string `json:"name"`
	}{
		Email:     email,
		Telephone: telephone,
		Name:      name,
	})
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
		return
	}
	url := strings.Join([]string{
		"http://",
		application.configs["app"].GetString("mailer_service"),
		"/user/",
		endpoint,
	}, "")
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
	}
}

func (h *httpRequests) mailerResetPasswordToken(token, email string) {
	body, err := json.Marshal(&struct {
		Token string `json:"token"`
		Email string `json:"email"`
	}{
		Token: token,
		Email: email,
	})
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
		return
	}
	url := strings.Join([]string{
		"http://",
		application.configs["app"].GetString("mailer_service"),
		"/user/pass/token",
	}, "")
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
	}
}

func (h *httpRequests) mailerResetPassword(email string) {
	body, err := json.Marshal(&struct {
		Email string `json:"email"`
	}{
		Email: email,
	})
	if err != nil {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
		return
	}
	url := strings.Join([]string{
		"http://",
		application.configs["app"].GetString("mailer_service"),
		"/user/pass/reset",
	}, "")
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println(runtimeinfo.Runtime(1), "; ERROR[", err, "]")
	}
}
