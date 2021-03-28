package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	episodeEndpoint = "https://api.spotify.com/v1/episodes/%s"
	trackEndpoint = "https://api.spotify.com/v1/tracks/%s"
)

func GetTrack(trackId, accessToken string) (*Track, error) {

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(trackEndpoint, trackId), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		resp := new(Track)
		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			return nil, err
		}
		return resp, nil
	}

	strResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("could not parse spotify track %s", string(strResp))
}

type Track struct {
	Album             Album        `json:"album"`
	Artists           []Artist     `json:"artists"`
	DurationMs        int          `json:"duration_ms"`
	Name              string       `json:"name"`
	URI               string       `json:"uri"`
}

type Artist struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
}

type Album struct {
	Images []Image `json:"images"`
}

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
