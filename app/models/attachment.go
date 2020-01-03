package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Attachment struct {
	ID                 uint         `gorm:"primary_key" json:"id"`

	Name               string       `json:"name"`

	UserId             uint         `json:"-"`
	User               User         `json:"-" gorm:"foreignkey:UserId"`

	MessageId          uint         `json:"-"`
	Message            Message      `json:"-" gorm:"foreignkey:MessageId"`

	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	DeletedAt          *time.Time   `json:"deleted_at"`
}

func (m *Attachment) BeforeCreate(scope *gorm.Scope) (err error) {

	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}

	return nil
}