package test

import (
	"dialogservice/mapper"
	"dialogservice/pckg/storage"
	"dialogservice/repository"
	"dialogservice/service"
	"log"
	"testing"
)

var (
	postgresOrm = storage.CreateConnection(
		storage.DBPostgres,
		storage.DSNPostgres,
		nil,
		"postgres",
		"toster123",
		"pet_finder_dialog",
		"5432",
		"disable",
	)
	repoApi    repository.DialogRepositoryAPI = repository.NewGormDialogRepositoryAPI(postgresOrm.DB)
	serviceApi                                = new(service.DialogServiceAPI)
)

func TestDownloadDialogsFlow(t *testing.T) {
	user := &mapper.UserViewModel{
		UserID: 1,
		Name:   "Vlad",
	}
	//
	view, err := serviceApi.DownloadDialogs(user, repoApi, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(view.Dialogs[0].Messages[0].MessageID)
	log.Println(view.Dialogs[0].Messages[len(view.Dialogs[0].Messages)-1].MessageID)

	//
	viewBatch, err := serviceApi.DownloadNextMessagesBatch(&mapper.NextMessagesBatchViewModel{
		DialogID:     view.Dialogs[0].DialogID,
		LastSkip:     view.Dialogs[0].SkipMessages,
		UserReceiver: user,
	}, repoApi, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(viewBatch.Messages[0].MessageID)
	log.Println(viewBatch.Messages[len(viewBatch.Messages)-1].MessageID)

	//
	viewBatch, err = serviceApi.DownloadNextMessagesBatch(&mapper.NextMessagesBatchViewModel{
		DialogID:     view.Dialogs[0].DialogID,
		LastSkip:     viewBatch.NextSkip,
		UserReceiver: user,
	}, repoApi, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(viewBatch.Messages[0].MessageID)
	log.Println(viewBatch.Messages[len(viewBatch.Messages)-1].MessageID)

	//
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddNewMessageFlow(t *testing.T) {
	user := &mapper.UserViewModel{
		UserID: 1,
		Name:   "Vlad",
	}
	//
	view, err := serviceApi.AddNewMessage(
		&mapper.AddNewMessageViewModel{
			DialogID:     1,
			Text:         "Hello",
			UserReceiver: user,
		}, repoApi, nil)
	if err != nil {
		t.Fatal(err)
	}
	if view.DialogID != 1 {
		t.Fatal("DialogID conflict")
	}
	log.Println(view)
	//
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateNewDialog(t *testing.T) {
	owner := &mapper.UserViewModel{
		UserID: 1,
		Name:   "Vlad",
	}
	//
	receiver := &mapper.UserViewModel{
		UserID: 4,
		Name:   "Ilya",
	}
	//
	view, err := serviceApi.CreateNewDialog(
		owner,
		receiver,
		repoApi,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(view)
	//
	viewReceiver, err := serviceApi.DownloadDialogs(receiver, repoApi, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(viewReceiver)
	//
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}
