package entity

import "time"

// Dialog : хранит всю информацию диалога.
//
type DialogEntity struct {
	//
	// PROPS
	//
	DialogID   uint64    `gorm:"primary_key;auto_increment"`
	DateCreate time.Time `gorm:"default:CURRENT_TIMESTAMP;not null;"`
	//
	// RELATED
	//
	//DialogUserList []DialogUserEntity //`gorm:"foreignKey:ForeignDialogID;references:DialogID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	//MessageList    []MessageEntity    //`gorm:"foreignKey:ForeignDialogID;references:DialogID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	//
	// gorm model
	//
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
}

// UserDialog : хранит всю информацию участника диалога.
//
type DialogUserEntity struct {
	//
	// FOREIGN KEY'S
	//
	//ForeignDialogID uint64 `gorm:"not null;foreignKey:ForeignDialogID;references:DialogID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	ForeignDialogID uint64
	Dialog          DialogEntity `gorm:"foreignKey:ForeignDialogID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//
	// PROPS
	//
	DialogName string `gorm:"uniqueIndex:member;size:255;not null"`
	//
	UserID       uint64 //`gorm:"primary_key;autoIncrement:false"`
	DialogUserID uint64 `gorm:"primary_key;auto_increment"`
	//
	UserName       string    `gorm:"uniqueIndex:member;size:255;not null"`
	ActivityStatus uint64    `gorm:"size:10;not null;"`
	DateCreate     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null;"`

	//
	// gorm model
	//
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
}

// Message : хранит всю информацию сообщения в диалоге.
//
type MessageEntity struct {
	//
	// FOREIGN KEY'S
	//
	ForeignDialogID uint64
	Dialog          DialogEntity `gorm:"foreignKey:ForeignDialogID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//
	// PROPS
	//
	MessageID  uint64    `gorm:"primary_key;auto_increment"`
	UserID     uint64    `gorm:"not null"`
	UserName   string    `gorm:"size:255;not null"`
	Text       string    `gorm:"size:255;"`
	DateCreate time.Time `gorm:"default:CURRENT_TIMESTAMP;not null;"`
	//
	// gorm model
	//
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
}
