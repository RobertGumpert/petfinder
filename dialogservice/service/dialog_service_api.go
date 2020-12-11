package service

import (
	"context"
	"dialogservice/entity"
	"dialogservice/mapper"
	"dialogservice/pckg/runtimeinfo"
	"dialogservice/repository"
	"log"
)

type DialogServiceAPI struct{}

func (s *DialogServiceAPI) CreateNewDialog(owner *mapper.UserViewModel, receiver *mapper.UserViewModel, db repository.DialogRepositoryAPI, ctx context.Context) (*mapper.CreateNewDialogViewModel, error) {
	if owner == nil || receiver == nil {
		return nil, mapper.ErrorNonValidData
	}
	if err := owner.Validator(); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	if err := receiver.Validator(); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	var (
		users = make([]*entity.UserEntity, 0)
	)
	users = append(users, owner.Mapper(), receiver.Mapper())
	id, err := db.CreateNewDialog(users, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.CreateNewDialogViewModel).Mapper(id), nil
}

func (s *DialogServiceAPI) DownloadDialogs(owner *mapper.UserViewModel, db repository.DialogRepositoryAPI, ctx context.Context) (*mapper.DownloadDialogsViewModel, error) {
	if err := owner.Validator(); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	messages, dialogs, dialogsUser, skip, err := db.DownloadDialogs(owner.UserID, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	viewModel := new(mapper.DownloadDialogsViewModel)
	if _, err := viewModel.Mapper(dialogs, dialogsUser, messages, owner.UserID, skip); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return viewModel, nil
}

func (s *DialogServiceAPI) DownloadNextMessagesBatch(viewModel *mapper.NextMessagesBatchViewModel, db repository.DialogRepositoryAPI, ctx context.Context) (*mapper.NextMessagesBatchResponse, error) {
	if err := viewModel.Validator(); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	messages, nextSkip, err := db.DownloadNextMessagesBatch(viewModel.DialogID, viewModel.LastSkip, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	var messagesViewModel []*mapper.MessageViewModel
	if messagesViewModel, err = new(mapper.MessageViewModel).MapperList(messages); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.NextMessagesBatchResponse).Mapper(viewModel.DialogID, nextSkip, viewModel.UserReceiver.UserID, messagesViewModel), nil
}

func (s *DialogServiceAPI) AddNewMessage(viewModel *mapper.AddNewMessageViewModel, db repository.DialogRepositoryAPI, ctx context.Context) (*mapper.AddNewMessageResponse, error) {
	if err := viewModel.Validator(); err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, err
	}
	messageEntity := viewModel.Mapper()
	_, err := db.AddNewMessage(messageEntity, ctx)
	if err != nil {
		go log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
		return nil, mapper.ErrorBadDataOperation
	}
	return new(mapper.AddNewMessageResponse).Mapper(messageEntity), nil
}
