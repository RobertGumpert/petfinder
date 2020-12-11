package test

import (
	"dialogservice/entity"
	"dialogservice/mapper"
	"dialogservice/pckg/storage"
	"dialogservice/repository"
	"log"
	"strconv"
	"testing"
	"time"
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

	api repository.DialogRepositoryAPI = repository.NewGormDialogRepositoryAPI(postgresOrm.DB)
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
		_, _ = api.AddNewMessage(&entity.MessageEntity{
			ForeignDialogID: 1,
			UserID:          2,
			UserName:        "Danil",
			Text:            strconv.Itoa(i),
		}, nil)
		time.Sleep(time.Second)
		_, _ = api.AddNewMessage(&entity.MessageEntity{
			ForeignDialogID: 2,
			UserID:          1,
			UserName:        "Vlad",
			Text:            "bbbb",
		}, nil)
		_, _ = api.AddNewMessage(&entity.MessageEntity{
			ForeignDialogID: 1,
			UserID:          1,
			UserName:        "Vlad",
			Text:            strconv.Itoa(i * 10),
		}, nil)
	}
}

func TestDownloadDialogs(t *testing.T) {

	// select * from message_entities where foreign_dialog_id = 1  ORDER BY date_create DESC
	// OFFSET 0
	// ROWS LIMIT 15;

	messages, dialogs, dialogsUser, skip, err := api.DownloadDialogs(1, nil)
	if err != nil {
		t.Fatal(err)
	}
	var viewModel = new(mapper.DownloadDialogsViewModel)
	if _, err := viewModel.Mapper(dialogs, dialogsUser, messages, 1, skip); err != nil {
		log.Println(err)
	}
	var mes  = make([]entity.MessageEntity, 0)
	for _, m := range messages {
		if m.ForeignDialogID == 1{
			mes = append(mes, m)
		}
	}
	log.Println("Skip -> ", skip)
	log.Println(mes[0].MessageID)
	log.Println(mes[len(mes)-1].MessageID)

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
