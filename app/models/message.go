package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Message struct {
	ID                     uint         `gorm:"primary_key" json:"id"`

	Content                string       `json:"data"`

	SenderId               uint         `json:"-"`
	Sender                 User         `json:"user" gorm:"foreignkey:SenderId"`

	ReceiverId             uint         `json:"-"`
	Receiver               User         `json:"-" gorm:"foreignkey:ReceiverId"`

	Attachments            []Attachment `json:"attachments"`

	Edited                 bool         `json:"edited" gorm:"default:false"`
	Deleted                bool         `json:"deleted" gorm:"default:false"`

	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
	DeletedAt              *time.Time   `json:"deleted_at"`
}

func (m *Message) BeforeCreate(scope *gorm.Scope) (err error) {

	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}

	if err = scope.SetColumn("CreatedAt", time.Now()); err != nil {
		return err
	}

	if err = scope.SetColumn("UpdatedAt", time.Now()); err != nil {
		return err
	}

	return nil
}