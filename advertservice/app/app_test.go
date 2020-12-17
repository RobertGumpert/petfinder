package app

import (
	"advertservice/mapper"
	"advertservice/pckg/conf"
	"advertservice/pckg/storage"
	"advertservice/repository"
	"advertservice/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	fileUpload = "C:/PetFinderRepos/petfinder/fileservice/test_jpeg.jpg"
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
	testRoot        = "C:/PetFinderRepos/petfinder/advertservice"
	testConfigs     map[string]*viper.Viper
	testPostgresOrm *storage.Storage
	testRepository  repository.AdvertRepository
	testSearchModel repository.SearchModel
	testService     *service.AdvertService
	testApplication *Application
)

func initFields() {
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
		"pet_finder_advert",
		"5432",
		"disable",
	)
	testRepository = repository.NewGormAdvertRepository(testPostgresOrm.DB)
	testSearchModel = repository.NewGormSquareSearchModel(
		testPostgresOrm.DB,
		mapper.CompareAdvertTime,
		mapper.OneKilometerScale * float64(1000000),
	)
	testService = service.NewAdvertService(
		mapper.LifetimeOfFoundAnimalAdvert,
		mapper.LifetimeOfLostAnimalAdvert,
		mapper.CompareAdvertTime,
	)
	testApplication = newTestApp(
		testConfigs,
		testService,
		testRepository,
		testSearchModel,
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

func postToken(srv *httptest.Server, endPoint, token string) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", srv.URL, endPoint), nil)
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

func postJSON(srv *httptest.Server, endPoint string, body interface{}) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", srv.URL, endPoint), structToIO(body))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	res, err := client.Do(req)
	return res, err
}

func postFormDataJSONFile(srv *httptest.Server, endPoint, token string, body interface{}) (*http.Response, error) {
	var (
		requestBuffer bytes.Buffer
		requestWriter = multipart.NewWriter(&requestBuffer)
		err           error
		//
		fileWriter io.Writer
		fileReader io.Reader = func() io.Reader {
			file, err := os.Open(fileUpload)
			if err != nil {
				log.Fatal(err)
			}
			return file
		}()
		//
		jsonWriter io.Writer
		jsonReader io.Reader = structToIO(body)
	)
	//
	if fileWriter, err = requestWriter.CreateFormFile("file", fileUpload); err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(fileWriter, fileReader); err != nil {
		log.Fatal(err)
	}
	if jsonWriter, err = requestWriter.CreateFormField("json"); err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(jsonWriter, jsonReader); err != nil {
		log.Fatal(err)
	}
	err = requestWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
	//
	client := http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", srv.URL, endPoint),
		&requestBuffer,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", requestWriter.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

func postFormDataJSON(srv *httptest.Server, endPoint, token string, body interface{}) (*http.Response, error) {
	var (
		requestBuffer bytes.Buffer
		requestWriter = multipart.NewWriter(&requestBuffer)
		err           error
		//
		jsonWriter io.Writer
		jsonReader io.Reader = structToIO(body)
	)
	//
	if jsonWriter, err = requestWriter.CreateFormField("json"); err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(jsonWriter, jsonReader); err != nil {
		log.Fatal(err)
	}
	err = requestWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
	//
	client := http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", srv.URL, endPoint),
		&requestBuffer,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", requestWriter.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

func TestAddAdvertFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	res, err := postFormDataJSONFile(
		srv,
		"/api/advert/user/add",
		firstUserToken,
		&mapper.AdvertViewModel{
			AdType:       uint64(mapper.TypeLost),
			AnimalType:   "Собака",
			AnimalBreed:  "Овчарка",
			GeoLatitude:  50.0001,
			GeoLongitude: 50.0001,
			CommentText:  "Потерялась, помогите найти!",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status isn't OK")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	viewModel := new(mapper.AdvertViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(viewModel.AdID)
	//
	//
	//
	res, err = getToken(
		srv,
		"/api/advert/user/list",
		firstUserToken,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("status isn't OK")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	listViewModel := new(mapper.ListAdvertViewModel)
	err = json.Unmarshal(body, listViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println("Lost : ", listViewModel.Lost.List)
	log.Println("Lost, expire : ", listViewModel.Lost.Expire)
	log.Println("Found : ", listViewModel.Found.List)
	log.Println("Found, expire : ", listViewModel.Found.Expire)
	//
	//
	//
	res, err = postFormDataJSON(
		srv,
		"/api/advert/user/add",
		secondUserToken,
		&mapper.AdvertViewModel{
			AdType:       uint64(mapper.TypeFound),
			AnimalType:   "Собака",
			AnimalBreed:  "Овчарка",
			GeoLatitude:  50.0002,
			GeoLongitude: 50.0003,
			CommentText:  "Нашлась овчарка",
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
	viewModel = new(mapper.AdvertViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(viewModel.AdID)
	//
	//
	//
	res, err = getToken(
		srv,
		"/api/advert/user/list",
		secondUserToken,
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
	listViewModel = new(mapper.ListAdvertViewModel)
	err = json.Unmarshal(body, listViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println("Lost : ", listViewModel.Lost.List)
	log.Println("Lost, expire : ", listViewModel.Lost.Expire)
	log.Println("Found : ", listViewModel.Found.List)
	log.Println("Found, expire : ", listViewModel.Found.Expire)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCloseRefreshAdvertFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	res, err := postJSONToken(
		srv,
		"/api/advert/user/close",
		firstUserToken,
		&mapper.IdentifierAdvertViewModel{
			AdType: uint64(mapper.TypeLost),
			AdID:   6,
		},
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
	viewModel := new(mapper.UpdateLifetimeViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(*viewModel)
	//
	//
	//
	res, err = postToken(
		srv,
		"/api/advert/user/list",
		firstUserToken,
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
	listViewModel := new(mapper.ListAdvertViewModel)
	err = json.Unmarshal(body, listViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println("Lost : ", listViewModel.Lost.List)
	log.Println("Lost, expire : ", listViewModel.Lost.Expire)
	log.Println("Found : ", listViewModel.Found.List)
	log.Println("Found, expire : ", listViewModel.Found.Expire)
	//
	//
	//
	res, err = postJSONToken(
		srv,
		"/api/advert/user/refresh",
		firstUserToken,
		&mapper.IdentifierAdvertViewModel{
			AdType: uint64(mapper.TypeLost),
			AdID:   6,
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
	viewModel = new(mapper.UpdateLifetimeViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println(*viewModel)
	//
	//
	//
	res, err = postToken(
		srv,
		"/api/advert/user/list",
		firstUserToken,
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
	listViewModel = new(mapper.ListAdvertViewModel)
	err = json.Unmarshal(body, listViewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println("Lost : ", listViewModel.Lost.List)
	log.Println("Lost, expire : ", listViewModel.Lost.Expire)
	log.Println("Found : ", listViewModel.Found.List)
	log.Println("Found, expire : ", listViewModel.Found.Expire)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchInAreaFlow(t *testing.T) {
	initFields()
	srv := httptest.NewServer(testApplication.HttpAPI)
	defer srv.Close()
	//
	//
	//
	res, err := postJSON(
		srv,
		"/api/advert/get/in/area",
		&mapper.SearchInAreaViewModel{
			AdOwnerID:     76,
			OnlyNotClosed: true,
			GeoLongitude:  60.0,
			GeoLatitude:   60.0,
		},
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
	viewModel := new(mapper.ListAdvertViewModel)
	err = json.Unmarshal(body, viewModel)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//
	//
	//
	log.Println("Lost : ", viewModel.Lost.List)
	log.Println("Lost, expire : ", viewModel.Lost.Expire)
	log.Println("Found : ", viewModel.Found.List)
	log.Println("Found, expire : ", viewModel.Found.Expire)
	//
	//
	//
	err = testPostgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}
