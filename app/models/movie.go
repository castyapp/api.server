package models

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Movie struct {
	ID                 uint          `gorm:"primary_key" json:"id"`

	MovieUri           string        `json:"movie_uri"`

	Poster             string        `json:"poster"`
	Subtitles          []Subtitle    `json:"subtitles"`

	Size               uint          `json:"size"`
	Duration           time.Duration `json:"duration"`
	LastPlayedTime     time.Duration `json:"last_played_time"`

	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

type MovieUrl string

type OpenSearchXmlResponse struct {
	ShortName     string `xml:"ShortName"`
	Description   string `xml:"Description"`
	InputEncoding string `xml:"InputEncoding"`
	SearchForm    string `xml:"SearchForm"`
}

func (mu MovieUrl) getParsedUrl() *url.URL {
	uri, err := url.Parse(string(mu))
	if err != nil {
		return nil
	}
	return uri
}

func (mu MovieUrl) IsPirateBay() bool {

	if uri := mu.getParsedUrl(); uri != nil {

		response, err := http.Get(fmt.Sprintf("http://%s/opensearch.xml", uri.Host))
		if err != nil {
			return false
		}

		var result OpenSearchXmlResponse
		decoder := xml.NewDecoder(response.Body)
		if err := decoder.Decode(&result); err != nil {
			return false
		}

		if strings.ContainsAny(result.Description, "pirate bay") {
			return true
		}

		if strings.ContainsAny(result.ShortName, "pirate bay") {
			return true
		}

		return false
	}

	return false
}

func (m *Movie) BeforeCreate(scope *gorm.Scope) (err error) {

	if err = scope.SetColumn("ID", uuid.New().ID()); err != nil {
		return err
	}

	return nil
}