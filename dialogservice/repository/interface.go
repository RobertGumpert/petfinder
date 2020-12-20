package repository

import (
	"context"
	"dialogservice/entity"
	"dialogservice/mapper"
)

type DialogRepository interface {
	createDialog(dialogEntity *entity.DialogEntity, ctx context.Context) (uint64, error)
	getDialog(dialogId uint64, ctx context.Context) (*entity.DialogEntity, error)
}

type MessageRepository interface {
	createMessage(message *entity.MessageEntity, ctx context.Context) (uint64, error)
	getMessageList(dialogId uint64, ctx context.Context) ([]entity.MessageEntity, error)
	getMessageByID(messageId uint64, ctx context.Context) (*entity.MessageEntity, error)
}

type DialogUserRepository interface {
	createDialogUser(dialogEntity *entity.DialogUserEntity, ctx context.Context) (uint64, error)
	getListByUser(userId uint64, ctx context.Context) ([]entity.DialogUserEntity, error)
	updateByUser(userId uint64, status mapper.DialogUserStatus, ctx context.Context)
	updateByIDs(ids []uint64, fields map[string]interface{}, ctx context.Context) error
}

type DialogRepositoryAPI interface {
	DialogUserRepository
	MessageRepository
	DialogRepository
	//
	CreateNewDialog(users []*entity.UserEntity, ctx context.Context) (uint64, error)
	AddNewMessage(message *entity.MessageEntity, ctx context.Context) (uint64, error)
	UpdateUserName(userId uint64, userName string, ctx context.Context) error
	DownloadDialogs(userId uint64, ctx context.Context) ([]entity.MessageEntity, []entity.DialogEntity, []entity.DialogUserEntity, uint64, error)
	DownloadNextMessagesBatch(dialogId uint64, lastSkip uint64, ctx context.Context) ([]entity.MessageEntity,  uint64, error)
}
