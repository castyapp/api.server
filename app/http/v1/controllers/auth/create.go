package auth

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/components/recaptcha"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/config"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// Create user jwt token
func Create(ctx *gin.Context) {

	var (
		request = &requests.CreateAuthTokenRequest{
			User: ctx.PostForm("user"),
			Pass: ctx.PostForm("pass"),
		}
		//rules = govalidator.MapData{
		//"pass": []string{"required", "min:4", "max:50"},
		//"user": []string{"required"},
		//}
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	if config.Map.Recaptcha.Enabled {
		if _, err := recaptcha.Verify(ctx); err != nil {
			sentry.CaptureException(err)
			ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
				"recaptcha": []string{
					"Captcha is invalid!",
				},
			}))
			return
		}
	}

	response, err := grpc.AuthServiceClient.Authenticate(ctx, &proto.AuthRequest{
		User: ctx.PostForm("user"),
		Pass: ctx.PostForm("pass"),
	})

	if err != nil {
		sentry.CaptureException(err)
		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":         string(response.Token),
		"refresh_token": string(response.RefreshedToken),
		"type":          "bearer",
	}))
}
