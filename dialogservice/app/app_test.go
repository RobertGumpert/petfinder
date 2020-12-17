package app

import (
	"bytes"
	"dialogservice/mapper"
	"dialogservice/pckg/conf"
	"dialogservice/pckg/storage"
	"dialogservice/repository"
	"dialogservice/service"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	secondUser = &mapper.UserViewModel{
		UserID: 77,
		Name:   "Danil",
	}
	secondUserToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05OTktOTk5LTk5OTkiLCJzZWNvbmQiOiJEYW5pbCJ9LCJleHAiOjE4Mjg5NDg3OTN9.jGU42jklINwz-W2jXyOwxjSqcnw1z-ygOIHS16uQsgo"
	//
	firstUser = &mapper.UserViewModel{
		UserID: 76,
		Name:   "Vlad",
	}
	firstUserToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTgyODk0ODcwNn0.e3REagO26al63HRHe75p2IokwihLGNW7xu284I5xhbE"
	//
	testRoot                = "C:/PetFinderRepos/petfinder/dialogservice"
	testConfigs             map[string]*viper.Viper
	testPostgresOrm         *storage.Storage
	testRepositoryDialogAPI repository.DialogRepositoryAPI
	testServiceDialogAPI    *service.DialogServiceAPI
	testApplication         *Application
)

func initFields() {
	testRoot = "C:/PetFinderRepos/petfinder/dialogservice"
	testConfigs = conf.ReadConfigs(
		testRoot,
		"app",
	)
	testPostgresOrm = storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		"postgres",
		"toster123",
		"pet_finder_dialog",
		"5432",
		"disable",
	)
	testRepositoryDialogAPI = repository.NewGormDialogRepositoryAPI(testPostgresOrm.DB)
	testServiceDialogAPI = service.NewDialogServiceAPI()
	testApplication = newTestApp(
		testConfigs,
		testServiceDialogAPI,
		testRepositoryDialogAPI,
	)
}

type errorResponseViewModel struct {
	Error string `json:"error"`
}

func structToIO(body interface{}) io.Reader {
	bts, _ := json.Marshal(body)
	ioReader := bytes.NewReader(bts)
	return ioReader
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

func TestCreateNewDialogFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	res, err := postJSONToken(
		srv,
		"/api/user/dialog/create",
		firstUserToken,
		secondUser,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status is OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel := new(mapper.CreateNewDialogViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(res.StatusCode)
	log.Println(*viewModel)
	//
	//
	//
	res, err = postJSONToken(
		srv,
		"/api/user/dialog/create",
		secondUserToken,
		firstUser,
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
	duplicateViewModel := new(mapper.CreateNewDialogViewModel)
	err = json.Unmarshal(body, duplicateViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if duplicateViewModel.ID != viewModel.ID {
		t.Fatal("ID's isn't compare ")
	}
	//
	//
	//
	log.Println(res.StatusCode)
	log.Println(*duplicateViewModel)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddMessagesFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	for i := 0; i < 45; i++ {
		res, err := postJSONToken(
			srv,
			"/api/user/message/send",
			firstUserToken,
			&mapper.AddNewMessageViewModel{
				DialogID: 21,
				Text:     strconv.Itoa(i),
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatal("status is OK")
		}
		//
		time.Sleep(time.Second * 1)
		//
		res, err = postJSONToken(
			srv,
			"/api/user/message/send",
			secondUserToken,
			&mapper.AddNewMessageViewModel{
				DialogID: 21,
				Text:     strconv.Itoa(i * 10),
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatal("status is OK")
		}
	}
	//
	//
	//
	err := testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadDialogFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	var printMessages = func(messages []*mapper.MessageViewModel) {
		if len(messages) == 0 {
			log.Println("List messages is empty")
			return
		}
		for _, m := range messages {
			log.Println("Message ID : ", m.MessageID, "; Text : ", m.Text, "; Date create : ", m.DateCreate)
		}
	}
	//
	//
	//
	res, err := getToken(
		srv,
		"/api/user/dialog/get",
		firstUserToken,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status is OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel := new(mapper.DownloadDialogsViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _, dialog := range viewModel.Dialogs {
		log.Println("ID : ", dialog.DialogID)
		log.Println("Dialog name : ", dialog.DialogName)
		log.Println("Skip :", dialog.SkipMessages)
		printMessages(dialog.Messages)
	}
	//
	//
	//
	res, err = postJSONToken(
		srv,
		"/api/user/message/batching/next",
		firstUserToken,
		&mapper.NextMessagesBatchViewModel{
			DialogID:     21,
			LastSkip:     viewModel.Dialogs[0].SkipMessages,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status is OK")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	nextBatchViewModel := new(mapper.NextMessagesBatchResponse)
	err = json.Unmarshal(body, nextBatchViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Dialog ID : ", nextBatchViewModel.DialogID)
	log.Println("Skip :", nextBatchViewModel.NextSkip)
	printMessages(nextBatchViewModel.Messages)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}
