package discord

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Discord struct {
	ClientId      string
	ClientSecret  string
	RedirectUri   string
	ResponseType  string
	Scope         string
}

type VerifyResponse struct {
	AccessToken   string `json:"access_token"`
	Scope         string `json:"scope"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
}

type User struct {
	Id             string  `json:"id"`
	Username       string  `json:"username"`
	Verified       bool    `json:"verified"`
	Locale         string  `json:"locale"`
	MFAEnabled     bool    `json:"mfa_enabled"`
	Flags          int     `json:"flags"`
	Avatar         string  `json:"avatar"`
	Discriminator  string  `json:"discriminator"`
	Email          string  `json:"email"`
}

func NewDiscordOAuth() *Discord {
	return &Discord{
		ClientId: os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectUri: os.Getenv("DISCORD_REDIRECT_URI"),
		ResponseType: "code",
		Scope: "identify email",
	}
}

func (d *Discord) Redirect(c *gin.Context)  {
	c.Redirect(302, "https://discordapp.com/api/oauth2/authorize?client_id=" + d.ClientId + "&redirect_uri=" + d.RedirectUri + "&response_type=" + d.ResponseType + "&scope=" + d.Scope)
}

func (d *Discord) GetUser(accessToken string) (User, error) {

	req, err := http.NewRequest("GET", "https://discordapp.com/api/v6/users/@me", nil)
	req.Header.Add("Authorization", "Bearer " + accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var user User
	jsonErr := json.Unmarshal(body, &user)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	if resp.StatusCode == http.StatusOK {
		user.Avatar = "https://cdn.discordapp.com/avatars/" + user.Id + "/" + user.Avatar + ".jpg"
		return user, nil
	}

	return user, errors.New(string(body))
}

func (d *Discord) Verify(c *gin.Context) (VerifyResponse, error) {
	data := url.Values{}
	data.Set("code", c.Query("code"))
	data.Set("client_id", d.ClientId)
	data.Set("client_secret", d.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", d.Scope)
	data.Set("redirect_uri", d.RedirectUri)

	req, err := http.NewRequest("POST", "https://discordapp.com/api/v6/oauth2/token", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var responseBody VerifyResponse
	jsonErr := json.Unmarshal([]byte (body), &responseBody)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	if resp.StatusCode == http.StatusOK {
		return responseBody, nil
	}

	return responseBody, errors.New(string(body))
}