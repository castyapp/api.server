package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetEpisode(episodeID, accessToken string) (*Episode, error) {

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(episodeEndpoint, episodeID), nil)
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
		resp := new(Episode)
		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			return nil, err
		}
		return resp, nil
	}

	strResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("could not parse spotify episode %s", strResp)
}

type Episode struct {
	Show       Show    `json:"show"`
	DurationMs int     `json:"duration_ms"`
	Name       string  `json:"name"`
	Images     []Image `json:"images"`
	URI        string  `json:"uri"`
}

type Show struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
