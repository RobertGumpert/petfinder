package repository

import (
	"context"
	"dialogservice/entity"
	"dialogservice/mapper"
	"gorm.io/gorm"
	"time"
)

type DialogAPIGormRepository struct {
	*gorm.DB
}

func NewDialogAPIGormRepository(DB *gorm.DB) *DialogAPIGormRepository {
	return &DialogAPIGormRepository{DB: DB}
}

func (d *DialogAPIGormRepository) createDialogUser(dialogEntity *entity.DialogUserEntity, ctx context.Context) (uint64, error) {
	tx := d.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(dialogEntity).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return dialogEntity.DialogUserID, nil
}

func (d *DialogAPIGormRepository) getListByUser(userId uint64, ctx context.Context) ([]entity.DialogUserEntity, error) {
	var (
		userDialogs []entity.DialogUserEntity
	)
	if err := d.Where(&entity.DialogUserEntity{UserID: userId}).Find(&userDialogs).Error; err != nil {
		return nil, err
	}
	return userDialogs, nil
}

func (d *DialogAPIGormRepository) updateByUser(userId uint64, status mapper.DialogUserStatus, ctx context.Context) {
	return
}

func (d *DialogAPIGormRepository) createMessage(message *entity.MessageEntity, ctx context.Context) (uint64, error) {
	tx := d.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return message.MessageID, tx.Commit().Error
}

func (d *DialogAPIGormRepository) getMessageList(dialogId uint64, ctx context.Context) ([]entity.MessageEntity, error) {
	return nil, nil
}

func (d *DialogAPIGormRepository) getMessageByID(messageId uint64, ctx context.Context) (*entity.MessageEntity, error) {
	return &entity.MessageEntity{}, nil
}

func (d *DialogAPIGormRepository) createDialog(dialogEntity *entity.DialogEntity, ctx context.Context) (uint64, error) {
	tx := d.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(dialogEntity).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return dialogEntity.DialogID, tx.Commit().Error
}

func (d *DialogAPIGormRepository) getDialog(dialogId uint64, ctx context.Context) (*entity.DialogEntity, error) {
	return &entity.DialogEntity{}, nil
}

//
// API
//

func (d *DialogAPIGormRepository) CreateNewDialog(users []*entity.UserEntity, ctx context.Context) (uint64, error) {
	dateCreate := time.Now()
	newDialog := new(entity.DialogEntity)
	newDialog.DateCreate = dateCreate
	id, err := d.createDialog(newDialog, ctx)
	if err != nil {
		return 0, err
	}
	for _, user := range users {
		_, err := d.createDialogUser(&entity.DialogUserEntity{
			ForeignDialogID: id,
			Dialog:          *newDialog,
			DateCreate:      dateCreate,
			UserID:          user.ID,
			UserName:        user.Name,
			ActivityStatus:  uint64(mapper.DialogUserStatusActive),
		}, nil)
		if err != nil {
			return 0, err
		}
	}
	return id, err
}

func (d *DialogAPIGormRepository) AddNewMessage(message *entity.MessageEntity, ctx context.Context) (uint64, error) {
	dateCreate := time.Now()
	message.DateCreate = dateCreate
	message.Dialog = entity.DialogEntity{
		DialogID: message.ForeignDialogID,
	}
	id, err := d.createMessage(message, ctx)
	return id, err
}

func (d *DialogAPIGormRepository) UpdateUserName(userId uint64, userName string, ctx context.Context) error {
	return nil
}

func (d *DialogAPIGormRepository) DownloadDialogs(userId uint64, ctx context.Context) ([]entity.MessageEntity, []entity.DialogEntity, uint64, error) {
	userDialogs, err := d.getListByUser(userId, ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	var (
		dialogs  []entity.DialogEntity
		messages []entity.MessageEntity
		ids      = make([]uint64, 0)
	)
	for _, userDialog := range userDialogs {
		dialogs = append(dialogs, entity.DialogEntity{
			DialogID: userDialog.ForeignDialogID,
		})
		ids = append(ids, userDialog.ForeignDialogID)
	}
	if err := d.Where("foreign_dialog_id IN ?", ids).Order("date_create desc").Limit(15).Find(&messages).Error; err != nil {
		return nil, nil, 0, err
	}
	return messages, dialogs, 0, nil
}

func (d *DialogAPIGormRepository) DownloadNextMessagesBatch(dialogId uint64, lastSkip uint64, ctx context.Context) ([]entity.MessageEntity, uint64, error) {
	var (
		messages []entity.MessageEntity
		nextSkip = lastSkip + 15
	)
	if err := d.Where("foreign_dialog_id = ?", dialogId).Order("date_create desc").Offset(int(nextSkip)).Limit(15).Find(&messages).Error; err != nil {
		return nil, nextSkip, err
	}
	return messages, nextSkip, nil
}
