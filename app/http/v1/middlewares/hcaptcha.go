package middlewares

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/castyapp/api.server/config"
	"github.com/MrJoshLab/go-respond"
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

func HcaptchaMiddleware(ctx *gin.Context) {

	hcaptchaHeader := ctx.GetHeader("h-captcha-response")
	if hcaptchaHeader == "" {
		ctx.AbortWithStatusJSON(respond.Default.ValidationErrors(map[string]interface{}{
			"recaptcha": []string{
				"Captcha is required!",
			},
		}))
		return
	}

	var (
		params                       = url.Values{}
		result                       = new(SiteVerificationResponse)
		token                        = ctx.GetHeader("h-captcha-response")
		invalidCode, invalidResponse = respond.Default.ValidationErrors(map[string]interface{}{
			"recaptcha": []string{
				"Captcha is invalid!",
			},
		})
	)

	params.Set("secret", config.Map.Recaptcha.Secret)
	params.Set("response", token)
	body := strings.NewReader(params.Encode())

	response, err := http.Post(siteVerifyEndpoint, "application/x-www-form-urlencoded", body)
	if err != nil {
		ctx.AbortWithStatusJSON(invalidCode, invalidResponse)
		return
	}

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		ctx.AbortWithStatusJSON(invalidCode, invalidResponse)
		return
	}

	if result.Success {
		ctx.Next()
		return
	}

	ctx.AbortWithStatusJSON(invalidCode, invalidResponse)
	return
}
