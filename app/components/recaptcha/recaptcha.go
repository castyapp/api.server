package recaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/castyapp/api.server/config"
	"github.com/gin-gonic/gin"
)

const (
	siteVerifyEndpoint = "https://hcaptcha.com/siteverify"
)

type SiteVerificationResponse struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Credit      bool     `json:"credit"`
}

func Verify(ctx *gin.Context) (*SiteVerificationResponse, error) {

	var (
		params = url.Values{}
		result = new(SiteVerificationResponse)
		token  = ctx.GetHeader("h-captcha-response")
	)

	params.Set("secret", config.Map.Recaptcha.Secret)
	params.Set("response", token)
	body := strings.NewReader(params.Encode())

	response, err := http.Post(siteVerifyEndpoint, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, err
	}

	if !result.Success {
		return result, fmt.Errorf("captcha is invalid! reason:[%v]", result.ErrorCodes)
	}

	return result, nil
}
