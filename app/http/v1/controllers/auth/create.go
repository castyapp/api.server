package auth

import (
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/recaptcha"
	"github.com/CastyLab/api.server/config"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// Create user jwt token
func Create(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"pass": []string{"required", "min:4", "max:50"},
			"user": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: false,
		}
		validator = govalidator.New(opts)
		validate  = validator.Validate()
	)

	if validate.Encode() == "" {

		if config.Map.App.Env == "prod" {
			if body, err := recaptcha.Verify(ctx); err != nil || !body.Success {
				sentry.CaptureException(err)
				ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
					"recaptcha": []string {
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

		ctx.JSON(respond.Default.Succeed(map[string] interface{} {
			"token": string(response.Token),
			"refreshed_token": string(response.RefreshedToken),
			"type": "bearer",
		}))
		return
	}

	validations := components.GetValidationErrorsFromGoValidator(validate)
	ctx.JSON(respond.Default.ValidationErrors(validations))
	return
}