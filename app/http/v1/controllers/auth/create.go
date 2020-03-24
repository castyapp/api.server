package auth

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/recaptcha"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"log"
	"net/http"
	"time"
)

// Create user jwt token
func Create(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"pass": []string{"required", "min:4", "max:30"},
			"user": []string{"required"},
			"g-recaptcha-response": []string{"required"},
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

		if success, err := recaptcha.Verify(ctx); err != nil || !success {
			ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
				"recaptcha": []string {
					"Captcha is invalid!",
				},
			}))
			return
		}

		userInput := ctx.PostForm("user")

		if userInput == "" {
			ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
				"user": []string {
					"User field is required!",
				},
			}))
			return
		}

		mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)
		response, err := grpc.AuthServiceClient.Authenticate(mCtx, &proto.AuthRequest{
			User: userInput,
			Pass: ctx.PostForm("pass"),
		})

		if err != nil {
			log.Println(err)
			return
		}

		if response.Code == http.StatusOK {
			ctx.JSON(respond.Default.Succeed(map[string] interface{} {
				"token": string(response.Token),
				"refreshed_token": string(response.RefreshedToken),
				"type": "bearer",
			}))
			return
		}

		ctx.JSON(respond.Default.SetStatusCode(http.StatusUnauthorized).
			SetStatusText("Failed!").
			RespondWithMessage("Unauthorized!"))
		return
	}

	validations := components.GetValidationErrorsFromGoValidator(validate)
	ctx.JSON(respond.Default.ValidationErrors(validations))
	return
}