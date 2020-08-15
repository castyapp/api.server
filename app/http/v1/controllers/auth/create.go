package auth

import (
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Create user jwt token
func Create(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"pass": []string{"required", "min:4", "max:50"},
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

		userInput := ctx.PostForm("user")

		if userInput == "" {
			ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
				"user": []string {
					"User field is required!",
				},
			}))
			return
		}

		md := metadata.New(map[string] string{
			"g-recaptcha-response": ctx.PostForm("g-recaptcha-response"),
		})

		response, err := grpc.AuthServiceClient.Authenticate(ctx, &proto.AuthRequest{
			User: userInput,
			Pass: ctx.PostForm("pass"),
		}, grpc2.Header(&md))

		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
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