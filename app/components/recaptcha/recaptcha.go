package recaptcha

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type VerifyResponse struct {
	Success       bool       `json:"success"`
	ChallengeTs   string     `json:"challenge_ts"`
	Hostname      string     `json:"hostname"`
	ErrorCodes    []string   `json:"error-codes"`
}

func Verify(ctx *gin.Context) (bool, error) {

	var (
		params     = url.Values{}
		remoteIP   = ctx.ClientIP()
		verifyResp = new(VerifyResponse)
		token      = ctx.PostForm("g-recaptcha-response")
	)

	params.Add("secret", os.Getenv("RECAPTCHA_SECRET_KEY"))
	params.Add("response", token)
	params.Add("remoteip", remoteIP)

	request, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(params.Encode()))
	if err != nil {
		return false, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&verifyResp); err != nil {
		return false, err
	}

	if verifyResp.Success == true {
		return true, nil
	}

	return false, errors.New("could not verify captcha")
}