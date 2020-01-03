package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Friend struct {
	ID            uint         `gorm:"primary_key" json:"id"`

	FriendId      uint         `json:"-"`
	Friend        User         `gorm:"foreignkey:FriendId" json:"friend"`

	UserId        uint         `json:"-"`
	User          User         `gorm:"foreignkey:UserId" json:"user"`

	Accepted      bool         `json:"accepted" gorm:"default:false"`

	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (f *Friend) BeforeCreate(scope *gorm.Scope) (err error) {
	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}
	return nil
}