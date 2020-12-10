package test

import (
	"dialogservice/entity"
	"dialogservice/pckg/storage"
	"dialogservice/repository"
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

	api repository.DialogRepositoryAPI = repository.NewDialogAPIGormRepository(postgresOrm.DB)
)

func foo() {
	if err := entity.GORMMigration(postgresOrm.DB); err != nil {
		log.Fatal(err)
	}
}

func TestCreateFlow(t *testing.T) {
	foo()
	//
	id, err := api.CreateNewDialog([]*entity.UserEntity{
		{
			ID:   1,
			Name: "Vlad",
		},
		{
			ID:   2,
			Name: "Danil",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(id)
	//
	id, err = api.CreateNewDialog([]*entity.UserEntity{
		{
			ID:   1,
			Name: "Vlad",
		},
		{
			ID:   3,
			Name: "Vika",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(id)
	//
	id, err = api.CreateNewDialog([]*entity.UserEntity{
		{
			ID:   3,
			Name: "Vika",
		},
		{
			ID:   2,
			Name: "Danil",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(id)
	//
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddMessage(t *testing.T)  {
	for i := 0; i < 30; i++ {
		_, _ = api.AddNewMessage(&entity.MessageEntity{
			ForeignDialogID: 2,
			UserID:          3,
			UserName:        "Vika",
			Text:            "aaa",
		}, nil)
	}
}

func TestDownloadDialogs(t *testing.T) {

	messages, _, skip, err := api.DownloadDialogs(1, nil)
	if err != nil {
		t.Fatal(err)
	}

	messages, skip, err = api.DownloadNextMessagesBatch(1, skip, nil)
	if len(messages) == 0{
		log.Println("len 0 ")
	} else {
		log.Println("Skip -> ", skip)
		log.Println(messages[0].MessageID)
		log.Println(messages[len(messages)-1].MessageID)
	}

	messages, skip, err = api.DownloadNextMessagesBatch(1, skip, nil)
	if len(messages) == 0{
		log.Println("len 0 ")
	} else {
		log.Println("Skip -> ", skip)
		log.Println(messages[0].MessageID)
		log.Println(messages[len(messages)-1].MessageID)
	}

	messages, skip, err = api.DownloadNextMessagesBatch(1, skip, nil)
	if len(messages) == 0{
		log.Println("len 0 ")
	} else {
		log.Println("Skip -> ", skip)
		log.Println(messages[0].MessageID)
		log.Println(messages[len(messages)-1].MessageID)
	}

	messages, skip, err = api.DownloadNextMessagesBatch(1, skip, nil)
	if len(messages) == 0{
		log.Println("len 0 ")
	} else {
		log.Println("Skip -> ", skip)
		log.Println(messages[0].MessageID)
		log.Println(messages[len(messages)-1].MessageID)
	}


	//
	err = postgresOrm.CloseConnection()
	if err != nil {
		t.Fatal(err)
	}
}
