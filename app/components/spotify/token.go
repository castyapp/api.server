package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	tokenEndpoint = "https://accounts.spotify.com/api/token"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int    `json:"expires_at"`
	Scope       string `json:"scope"`
}

func GetAuthenticationToken() (*Token, error) {

	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", os.Getenv("SPOTIFY_OAUTH_CLIENT_ID"))
	params.Set("client_secret", os.Getenv("SPOTIFY_OAUTH_CLIENT_SECRET"))

	request, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	resp := new(Token)
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
