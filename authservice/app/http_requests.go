package app

import (
	"authservice/pckg/runtimeinfo"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
