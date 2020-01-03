package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Role struct {
	ID            uint         `gorm:"primary_key" json:"id"`
	Title         string       `json:"title"`
	CreatedAt     time.Time    `json:"-"`
	UpdatedAt     time.Time    `json:"-"`
}

func (r *Role) IsAdmin() bool {
	if r.Title == "admin" {
		return true
	}
	return false
}

func (r *Role) IsUser() bool {
	if r.Title == "user" {
		return true
	}
	return false
}

func (r *Role) BeforeCreate(scope *gorm.Scope) (err error) {

	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}

	return nil
}