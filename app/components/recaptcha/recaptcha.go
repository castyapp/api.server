package recaptcha

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	siteVerifyEndpoint = "https://hcaptcha.com/siteverify"
)

type SiteVerificationResponse struct {
	Success     bool      `json:"success"`
	ErrorCodes  []string  `json:"error-codes"`
	ChallengeTs string    `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	Credit      bool      `json:"credit"`
}

func Verify(ctx *gin.Context) (*SiteVerificationResponse, error) {

	var (
		params = url.Values{}
		result = new(SiteVerificationResponse)
		token  = ctx.GetHeader("h-captcha-response")
	)

	params.Set("secret", os.Getenv("RECAPTCHA_SECRET_KEY"))
	params.Set("response", token)
	body := strings.NewReader(params.Encode())

	response, err := http.Post(siteVerifyEndpoint, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}