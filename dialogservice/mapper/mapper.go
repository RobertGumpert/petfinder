package mapper

import (
	"dialogservice/entity"
	"strings"
	"time"
)

type UserViewModel struct {
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func (m *UserViewModel) Validator() error {
	if m.UserID == 0 || strings.TrimSpace(m.Name) == "" {
		return ErrorNonValidData
	}
	return nil
}

func (m *UserViewModel) Mapper() *entity.UserEntity {
	u := new(entity.UserEntity)
	u.ID = m.UserID
	u.Name = m.Name
	return u
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type CreateNewDialogViewModel struct {
	ID uint64 `json:"id"`
}

func (m *CreateNewDialogViewModel) Mapper(id uint64) *CreateNewDialogViewModel {
	m.ID = id
	return m
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type AddNewMessageResponse struct {
	DialogID       uint64    `json:"dialog_id"`
	UserReceiverID uint64    `json:"user_receiver_id"`
	MessageID      uint64    `json:"message_id"`
	UserID         uint64    `json:"user_id"`
	UserName       string    `json:"user_name"`
	Text           string    `json:"text"`
	DateCreate     time.Time `json:"date_create"`
}

func (m *AddNewMessageResponse) Mapper(message *entity.MessageEntity) *AddNewMessageResponse {
	m.Text = message.Text
	m.UserReceiverID = message.UserID
	m.DialogID = message.ForeignDialogID
	m.UserName = message.UserName
	m.DateCreate = message.DateCreate
	m.MessageID = message.MessageID
	return m
}

type AddNewMessageViewModel struct {
	DialogID     uint64 `json:"dialog_id"`
	Text         string `json:"text"`
	UserReceiver *UserViewModel
}

func (m *AddNewMessageViewModel) Validator() error {
	if m.DialogID == 0 {
		return ErrorNonValidData
	}
	if m.UserReceiver == nil {
		return ErrorNonValidData
	}
	return m.UserReceiver.Validator()
}

func (m *AddNewMessageViewModel) Mapper() *entity.MessageEntity {
	mes := new(entity.MessageEntity)
	mes.ForeignDialogID = m.DialogID
	mes.UserID = m.UserReceiver.UserID
	mes.UserName = m.UserReceiver.Name
	mes.Text = m.Text
	return mes
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type MessageViewModel struct {
	MessageID  uint64    `json:"message_id"`
	DialogID   uint64    `json:"dialog_id"`
	UserID     uint64    `json:"user_id"`
	UserName   string    `json:"user_name"`
	Text       string    `json:"text"`
	DateCreate time.Time `json:"date_create"`
}

func (m *MessageViewModel) Mapper(message *entity.MessageEntity) (*MessageViewModel, error) {
	if message.UserID == 0 || message.MessageID == 0 || message.ForeignDialogID == 0 {
		return nil, ErrorNonValidData
	}
	if strings.TrimSpace(message.UserName) == "" {
		return nil, ErrorNonValidData
	}
	m.DialogID = message.ForeignDialogID
	m.MessageID = message.MessageID
	m.UserID = message.UserID
	m.Text = message.Text
	m.UserName = message.UserName
	m.DateCreate = message.DateCreate
	return m, nil
}

func (m *MessageViewModel) MapperList(messages []entity.MessageEntity) ([]*MessageViewModel, error) {
	messagesViewModel := make([]*MessageViewModel, 0)
	for _, message := range messages {
		next := new(MessageViewModel)
		if _, err := next.Mapper(&message); err != nil {
			return nil, err
		}
		messagesViewModel = append(
			messagesViewModel,
			next,
		)
	}
	return messagesViewModel, nil
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type DialogViewModel struct {
	DialogID       uint64              `json:"dialog_id"`
	DialogName     string              `json:"dialog_name"`
	UserReceiverID uint64              `json:"user_receiver_id"`
	SkipMessages   uint64              `json:"skip_messages"`
	Messages       []*MessageViewModel `json:"messages"`
}

func (m *DialogViewModel) Mapper(dialog *entity.DialogEntity, dialogUser *entity.DialogUserEntity, messages []entity.MessageEntity, receiverId, skip uint64) (*DialogViewModel, error) {
	if dialog.DialogID != dialogUser.ForeignDialogID {
		return nil, ErrorNonValidData
	}
	if dialogUser.UserID != receiverId {
		return nil, ErrorNonValidData
	}
	m.DialogID = dialog.DialogID
	m.UserReceiverID = receiverId
	m.DialogName = dialogUser.DialogName
	m.SkipMessages = skip
	m.Messages = make([]*MessageViewModel, 0)
	for _, message := range messages {
		if message.ForeignDialogID != m.DialogID {
			return nil, ErrorMessageDoesntBelongDialog
		}
		nextMessage, err := new(MessageViewModel).Mapper(&message)
		if err != nil {
			return nil, err
		}
		m.Messages = append(
			m.Messages,
			nextMessage,
		)
	}
	return m, nil
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type NextMessagesBatchResponse struct {
	DialogID       uint64              `json:"dialog_id"`
	NextSkip       uint64              `json:"next_skip"`
	UserReceiverID uint64              `json:"user_receiver_id"`
	Messages       []*MessageViewModel `json:"messages"`
}

func (m *NextMessagesBatchResponse) Mapper(dialogId, nextSkip, userReceiverID uint64, messages []*MessageViewModel) *NextMessagesBatchResponse {
	m.DialogID = dialogId
	m.NextSkip = nextSkip
	m.UserReceiverID = userReceiverID
	m.Messages = messages
	return m
}

type NextMessagesBatchViewModel struct {
	DialogID     uint64 `json:"dialog_id"`
	LastSkip     uint64 `json:"last_skip"`
	UserReceiver *UserViewModel
}

func (m *NextMessagesBatchViewModel) Validator() error {
	if m.DialogID == 0 {
		return ErrorNonValidData
	}
	if m.UserReceiver == nil {
		return ErrorNonValidData
	}
	return m.UserReceiver.Validator()
}

//
//----------------------------------------------------------------------------------------------------------------------
//

type DownloadDialogsViewModel struct {
	Dialogs []*DialogViewModel `json:"dialogs"`
}

func (m *DownloadDialogsViewModel) Mapper(dialogs []entity.DialogEntity, dialogsUser []entity.DialogUserEntity, messages []entity.MessageEntity, receiverId, skip uint64) (*DownloadDialogsViewModel, error) {
	m.Dialogs = make([]*DialogViewModel, 0)
	for _, dialog := range dialogs {
		var (
			dialogId = dialog.DialogID
		)
		for _, dialogsUser := range dialogsUser {
			if dialogsUser.ForeignDialogID == dialogId {
				var (
					dialogViewModel = new(DialogViewModel)
					mes             = make([]entity.MessageEntity, 0)
				)
				for _, message := range messages {
					if message.ForeignDialogID == dialogId {
						mes = append(mes, message)
					}
				}
				if _, err := dialogViewModel.Mapper(&dialog, &dialogsUser, mes, receiverId, skip); err != nil {
					return nil, err
				}
				m.Dialogs = append(m.Dialogs, dialogViewModel)
			}
		}
	}
	return m, nil
}
