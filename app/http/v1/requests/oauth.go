package requests

type OauthCallbackRequest struct {
	Code string `validate:"required" form:"code"`
}
