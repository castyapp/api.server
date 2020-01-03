package models

import (
	"time"
)

type Theater struct {
	ID                 uint         `gorm:"primary_key" json:"id"`
	Title              string       `json:"title"`

	// Theater hash. this will use for websocket connections
    Hash               string       `json:"hash"`

	// 0 is just for the user
	// 1 is for everyone
	// 2 is for friends
	Privacy            int          `json:"privacy" gorm:"default:1"`

	// 0 is just for the user
	// 1 is for everyone
	// 2 is for friends
	VideoPlayerAccess  int          `json:"video_player_access" gorm:"default:1"`

	UserId             uint         `json:"-"`
	User               User         `json:"-" gorm:"foreignkey:UserId"`

	MovieId            uint        `json:"-"`
	Movie              Movie       `json:"movie" gorm:"foreignkey:MovieId"`

	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}