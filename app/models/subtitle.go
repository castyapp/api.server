package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Subtitle struct {
	ID                 uint       `gorm:"primary_key" json:"id"`

	Lang               string     `json:"size"`
	File               string     `json:"file"`

	Movie              Movie      `json:"movie" gorm:"foreignkey:MovieId"`
	MovieId            uint       `json:"-"`

	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (s *Subtitle) BeforeCreate(scope *gorm.Scope) (err error) {

	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}

	return nil
}