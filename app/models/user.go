package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"movie-night/app/components/strings"
	"time"
)

type User struct {
	ID            uint         `gorm:"primary_key" json:"id"`

	Fullname      string       `json:"fullname"`
	Username      string       `json:"username"`

	Hash          string       `json:"hash"`

	Email         string       `json:"email"`
	Password      string       `json:"-"`

	IsActive      bool         `json:"is_active" gorm:"default:true"`
	State         int          `json:"state" gorm:"default:0"`
	Activity      int          `json:"activity" gorm:"default:0"`

	Avatar        string       `json:"avatar"`

	RoleId        uint         `json:"-"`
	Role          Role         `json:"-" gorm:"foreignkey:RoleId"`

	LastLogin     time.Time    `json:"last_login"`
	JoinedAt      time.Time    `json:"joined_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (user *User) IsAdmin() bool {
	return user.Role.IsAdmin()
}

func (user *User) IsUser() bool {
	return user.Role.IsUser()
}

func (user *User) SetPassword(password string) {
	user.Password = user.HashPassword(password)
}

func (User) HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) GenerateUserHash() {
	user.Hash = strings.Random(30)
}

func (user *User) BeforeCreate(scope *gorm.Scope) (err error) {

	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}

	user.GenerateUserHash()
	user.SetPassword(user.Password)

	if err = scope.SetColumn("Avatar", "default"); err != nil {
		return err
	}

	if user.Role.ID == 0 && user.RoleId == 0 {
		role := new(Role)
		err = scope.DB().Table("roles").First(role, map[string] interface{} {"title": "user"}).Error
		if err != nil {
			return err
		}
		user.Role = *role
	}

	if err = scope.SetColumn("JoinedAt", time.Now()); err != nil {
		return err
	}

	if err = scope.SetColumn("LastLogin", time.Now()); err != nil {
		return err
	}

	if err = scope.SetColumn("RoleId", user.Role.ID); err != nil {
		return err
	}

	return err
}