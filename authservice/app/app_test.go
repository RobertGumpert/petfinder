package app

import (
	"authservice/mapper"
	"authservice/pckg/conf"
	"authservice/pckg/storage"
	"authservice/repository"
	"authservice/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type errorResponseViewModel struct {
	Error string `json:"error"`
}

type resetPasswordNewAccess struct {
	Token string `json:"token"`
}

type authViewModel struct {
	Token string                `json:"token"`
	User  *mapper.UserViewModel `json:"user"`
}

var (
	testRoot           = "C:/PetFinderRepos/petfinder/authservice"
	testConfigs        map[string]*viper.Viper
	testPostgresOrm    *storage.Storage
	testUserRepository repository.UserRepository
	testUserService    *service.User
	testApplication    *Application
)

func initFields() {
	testRoot = "C:/PetFinderRepos/petfinder/authservice"
	testConfigs = conf.ReadConfigs(
		testRoot,
		"app",
		"event_receivers",
	)
	testPostgresOrm = storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		"postgres",
		"toster123",
		"pet_finder_user",
		"5432",
		"disable",
	)
	testUserRepository = repository.NewUserGormRepository(testPostgresOrm.DB)
	testUserService = service.NewUserService(
		[]byte("auth_service"),
		5*time.Second,
		30*time.Hour,
		30*time.Minute,
	)
	testApplication = newTestApp(
		testConfigs,
		testUserService,
		testUserRepository,
	)
}

func structToIO(body interface{}) io.Reader {
	bts, _ := json.Marshal(body)
	ioReader := bytes.NewReader(bts)
	return ioReader
}

func postJSON(srv *httptest.Server, endPoint string, body interface{}) (*http.Response, error) {
	res, err := http.Post(fmt.Sprintf("%s%s", srv.URL, endPoint), "application/json", structToIO(body))
	return res, err
}

func getToken(srv *httptest.Server, endPoint, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", srv.URL, endPoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := http.Client{}
	res, err := client.Do(req)
	return res, err
}

func postJSONToken(srv *httptest.Server, endPoint, token string, body interface{}) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", srv.URL, endPoint), structToIO(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := http.Client{}
	res, err := client.Do(req)
	return res, err
}

func TestRegisterFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	newUser := &mapper.RegisterUserViewModel{
		Telephone: "8-000-000-0000",
		Password:  "test",
		Email:     "test@mail.ru",
		Name:      "test",
	}
	res, err := postJSON(srv, "/api/user/register", newUser)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	resNewUser := new(mapper.UserViewModel)
	err = json.Unmarshal(body, resNewUser)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(res.Status)
	log.Println(resNewUser)
	//
	//
	//
	res, err = postJSON(srv, "/api/user/register", newUser)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == http.StatusOK {
		t.Errorf("status is OK")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	errViewModel := new(errorResponseViewModel)
	err = json.Unmarshal(body, errViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(res.Status)
	log.Println(errViewModel)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthorizedFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	authUser := &mapper.AuthorizationUserViewModel{
		UserID:    74,
		Telephone: "8-000-000-0000",
		Password:  "test",
		Email:     "test@mail.ru",
	}
	res, err := postJSON(srv, "/api/user/authorized", authUser)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status not OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel := new(authViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	accessToken := viewModel.Token
	//
	//
	//
	log.Println(res.Status)
	log.Println(*viewModel.User)
	log.Println(viewModel.Token)
	//
	//
	//
	res, err = postJSON(srv, "/api/user/authorized", authUser)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == http.StatusOK {
		t.Fatal("status is OK")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	errViewModel := new(errorResponseViewModel)
	err = json.Unmarshal(body, errViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(res.Status)
	log.Println(errViewModel)
	log.Println("Time out")
	time.Sleep(time.Second * 6)
	//
	//
	//
	res, err = getToken(
		srv,
		"/api/user/access",
		accessToken,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == http.StatusOK {
		t.Fatal("status is OK")
	}
	//
	//
	//
	log.Println(res.Status)
	//
	//
	//
	res, err = getToken(
		srv,
		"/api/user/access/update",
		accessToken,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status not OK")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel = new(authViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(res.Status)
	log.Println(*viewModel.User)
	log.Println(viewModel.Token)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateUserFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	updateUser := &mapper.UpdateUserViewModel{
		UserID: 69,
		Name:   "Test Vlad",
	}
	oldAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U"
	res, err := postJSONToken(
		srv,
		"/api/user/update",
		oldAccessToken,
		updateUser,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status not OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel := new(mapper.UserViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	newAccessToken := viewModel.AccessToken
	if newAccessToken == oldAccessToken {
		t.Fatal("Access token is valid. ")
	}
	//
	//
	//
	log.Println(newAccessToken)
	//
	//
	//
	res, err = getToken(
		srv,
		"/api/user/access",
		newAccessToken,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status is OK")
	}
	//
	//
	//
	log.Println(res.Status)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetResetPassword(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	updateUser := &mapper.ResetPasswordViewModel{
		Email:     "walkmanmail19@gmail.com",
		Telephone: "8-999-999-9999",
	}
	oldAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U"
	res, err := postJSONToken(
		srv,
		"/api/user/password/token",
		oldAccessToken,
		updateUser,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status not OK")
	}
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestResetPassword(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	updateUser := &mapper.ResetPasswordViewModel{
		ResetToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiT0MwNU9Ua3RPVGs1TFRrNU9Uaz0iLCJzZWNvbmQiOiJ0b3N0ZXIxMjMifSwiZXhwIjoxNjA4MDU5MTgxfQ.TpLGN890ULFMyaJdK2ic5WxTivBkHofcZp_Pnd2SCUM",
		Email:      "walkmanmail19@gmail.com",
		Telephone:  "8-999-999-9999",
		Password:   "toster1234567",
	}
	oldAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U"
	res, err := postJSONToken(
		srv,
		"/api/user/password/reset",
		oldAccessToken,
		updateUser,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status not OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel := new(resetPasswordNewAccess)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(viewModel.Token)
	//
	//
	//
	res, err = getToken(
		srv,
		"/api/user/access",
		viewModel.Token,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status is OK")
	}
	//
	//
	//
	log.Println(res.Status)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}
