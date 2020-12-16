package app

import (
	"dialogservice/mapper"
	"dialogservice/pckg/runtimeinfo"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
	req, err := http.NewRequest("GET", url, nil)
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
