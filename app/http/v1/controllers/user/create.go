package user

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/components/recaptcha"
	"github.com/castyapp/api.server/config"
	"github.com/castyapp/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// Create a new user
func Create(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"fullname":              []string{"min:4", "max:30"},
			"password":              []string{"required", "min:4", "max:80"},
			"password_confirmation": []string{"required", "min:4", "max:80"},
			"username":              []string{"required", "between:3,20"},
			"email":                 []string{"required", "email"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		mCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	)

	defer cancel()

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	if ctx.PostForm("password") != ctx.PostForm("password_confirmation") {
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"password": []string{
				"Passwords are not match!",
			},
		}))
		return
	}

	if config.Map.Recaptcha.Enabled {
		if _, err := recaptcha.Verify(ctx); err != nil {
			ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
				"recaptcha": []string{
					"Captcha is invalid!",
				},
			}))
			return
		}
	}

	response, err := grpc.UserServiceClient.CreateUser(mCtx, &proto.CreateUserRequest{
		User: &proto.User{
			Fullname: ctx.PostForm("fullname"),
			Username: ctx.PostForm("username"),
			Email:    ctx.PostForm("email"),
			Password: ctx.PostForm("password"),
		},
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("failed").
			RespondWithMessage("Could not create user."))
		return
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":           string(response.Token),
		"refreshed_token": string(response.Token),
		"type":            "bearer",
	}))
	return
}
