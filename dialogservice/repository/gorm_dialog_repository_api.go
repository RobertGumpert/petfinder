package repository

import (
	"context"
	"dialogservice/entity"
	"dialogservice/mapper"
	"dialogservice/pckg/runtimeinfo"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type GormDialogRepositoryAPI struct {
	*gorm.DB
}

func NewGormDialogRepositoryAPI(DB *gorm.DB) *GormDialogRepositoryAPI {
	return &GormDialogRepositoryAPI{DB: DB}
}

func (d *GormDialogRepositoryAPI) createDialogUser(dialogEntity *entity.DialogUserEntity, ctx context.Context) (uint64, error) {
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

func (d *GormDialogRepositoryAPI) getListByUser(userId uint64, ctx context.Context) ([]entity.DialogUserEntity, error) {
	var (
		userDialogs []entity.DialogUserEntity
	)
	if err := d.Where(&entity.DialogUserEntity{UserID: userId}).Find(&userDialogs).Error; err != nil {
		return nil, err
	}
	return userDialogs, nil
}

func (d *GormDialogRepositoryAPI) updateByUser(userId uint64, status mapper.DialogUserStatus, ctx context.Context) {
	return
}

func (d *GormDialogRepositoryAPI) createMessage(message *entity.MessageEntity, ctx context.Context) (uint64, error) {
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

func (d *GormDialogRepositoryAPI) getMessageList(dialogId uint64, ctx context.Context) ([]entity.MessageEntity, error) {
	return nil, nil
}

func (d *GormDialogRepositoryAPI) getMessageByID(messageId uint64, ctx context.Context) (*entity.MessageEntity, error) {
	return &entity.MessageEntity{}, nil
}

func (d *GormDialogRepositoryAPI) createDialog(dialogEntity *entity.DialogEntity, ctx context.Context) (uint64, error) {
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

func (d *GormDialogRepositoryAPI) getDialog(dialogId uint64, ctx context.Context) (*entity.DialogEntity, error) {
	return &entity.DialogEntity{}, nil
}

func (d *GormDialogRepositoryAPI) updateByIDs(ids []uint64, fields map[string]interface{}, ctx context.Context) error {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entity.DialogEntity{}).Where("ad_id IN ? ", ids).Updates(fields).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//
// API
//

func (d *GormDialogRepositoryAPI) CreateNewDialog(users []*entity.UserEntity, ctx context.Context) (uint64, error) {
	var (
		addDialogUser = func(ownerID, receiverID uint64, dateCreate time.Time, newDialog entity.DialogEntity) error {
			_, err := d.createDialogUser(&entity.DialogUserEntity{
				ForeignDialogID: newDialog.DialogID,
				Dialog:          newDialog,
				DateCreate:      dateCreate,
				UserID:          users[ownerID].ID,
				UserName:        users[ownerID].Name,
				DialogName:      users[receiverID].Name,
				ActivityStatus:  uint64(mapper.DialogUserStatusActive),
			}, nil)
			return err
		}
		existRowByUserName entity.DialogUserEntity
	)
	err := d.Where(&entity.DialogUserEntity{UserName: users[0].Name, DialogName: users[1].Name}).First(&existRowByUserName).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == nil {
		return existRowByUserName.ForeignDialogID, nil
	}
	dateCreate := time.Now()
	newDialog := new(entity.DialogEntity)
	newDialog.DateCreate = dateCreate
	dialogId, err := d.createDialog(newDialog, ctx)
	if err != nil {
		return 0, err
	}
	newDialog.DialogID = dialogId
	if err := addDialogUser(0, 1, dateCreate, *newDialog); err != nil {
		return 0, err
	}
	if err := addDialogUser(1, 0, dateCreate, *newDialog); err != nil {
		return 0, err
	}
	return dialogId, err
}

func (d *GormDialogRepositoryAPI) AddNewMessage(message *entity.MessageEntity, ctx context.Context) (uint64, error) {
	dateCreate := time.Now()
	message.DateCreate = dateCreate
	message.Dialog = entity.DialogEntity{
		DialogID: message.ForeignDialogID,
	}
	id, err := d.createMessage(message, ctx)
	return id, err
}

func (d *GormDialogRepositoryAPI) UpdateUserName(userId uint64, userName string, ctx context.Context) error {
	dialogsUser, err := d.getListByUser(userId, ctx)
	if err != nil {
		return err
	}
	ids := make([]uint64, 0)
	for _, userDialog := range dialogsUser {
		ids = append(ids, userDialog.DialogUserID)
	}
	return d.updateByIDs(ids, map[string]interface{}{
		"user_name": userName,
	}, ctx)
}

func (d *GormDialogRepositoryAPI) DownloadDialogs(userId uint64, ctx context.Context) ([]entity.MessageEntity, []entity.DialogEntity, []entity.DialogUserEntity, uint64, error) {
	dialogsUser, err := d.getListByUser(userId, ctx)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	var (
		dialogs   []entity.DialogEntity
		messages  = make([]entity.MessageEntity, 0)
		dialogIds = make([]uint64, 0)
		wg        sync.WaitGroup
		skip      uint64 = 0
	)
	for _, userDialog := range dialogsUser {
		dialogs = append(dialogs, entity.DialogEntity{
			DialogID: userDialog.ForeignDialogID,
		})
		dialogIds = append(dialogIds, userDialog.ForeignDialogID)
	}
	for _, dialogId := range dialogIds {
		wg.Add(1)
		go func(wg *sync.WaitGroup, dialogId uint64) {
			defer wg.Done()
			var dialogMessages []entity.MessageEntity
			if err := d.Where("foreign_dialog_id = ?", dialogId).Order("date_create desc").Limit(mapper.BatchSizeMessages).Find(&dialogMessages).Error; err != nil {
				log.Println(runtimeinfo.Runtime(1), "; ERROR=[", err, "]")
				return
			}
			messages = append(messages, dialogMessages...)
		}(&wg, dialogId)
	}
	wg.Wait()
	return messages, dialogs, dialogsUser, skip, nil
}

func (d *GormDialogRepositoryAPI) DownloadNextMessagesBatch(dialogId uint64, lastSkip uint64, ctx context.Context) ([]entity.MessageEntity, uint64, error) {
	var (
		messages []entity.MessageEntity
		nextSkip = lastSkip + uint64(mapper.BatchSizeMessages)
	)
	if err := d.Where("foreign_dialog_id = ?", dialogId).Order("date_create desc").Offset(int(nextSkip)).Limit(mapper.BatchSizeMessages).Find(&messages).Error; err != nil {
		return nil, nextSkip, err
	}
	return messages, nextSkip, nil
}
