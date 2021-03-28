package user

import (
	"context"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/components/recaptcha"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/config"
	"github.com/castyapp/api.server/grpc"
	"github.com/gin-gonic/gin"
)

// Create a new user
func Create(ctx *gin.Context) {

	var (
		mCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
		request      = &requests.CreateUserRequest{
			Fullname:             ctx.PostForm("username"),
			Username:             ctx.PostForm("username"),
			Email:                ctx.PostForm("email"),
			Password:             ctx.PostForm("password"),
			PasswordConfirmation: ctx.PostForm("password_confirmation"),
		}
	)

	defer cancel()

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	if request.Password != request.PasswordConfirmation {
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
			Fullname: request.Fullname,
			Username: request.Username,
			Email:    request.Email,
			Password: request.Password,
		},
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":         string(response.Token),
		"refresh_token": string(response.Token),
		"type":          "bearer",
	}))
	return
}
