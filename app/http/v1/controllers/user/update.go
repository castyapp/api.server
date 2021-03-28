package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/components/strings"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/api.server/storage"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
)

func Update(ctx *gin.Context) {

	var (
		protoUser = new(proto.User)
		token     = ctx.GetHeader("Authorization")
		fullname  = ctx.PostForm("fullname")
	)

	if avatarFile, err := ctx.FormFile("avatar"); err == nil {
		avatarName := strings.RandomNumber(20)
		avatarObject := fmt.Sprintf("%s.png", avatarName)

		afile, err := avatarFile.Open()
		if err != nil {
			sentry.CaptureException(err)
			ctx.JSON(respond.Default.
				SetStatusText("Failed!").
				SetStatusCode(400).
				RespondWithMessage("Upload failed. Please try again later!"))
		}

		_, err = storage.Client.PutObjectWithContext(ctx, "avatars", avatarObject, afile, -1, minio.PutObjectOptions{})
		if err != nil {
			sentry.CaptureException(err)
			ctx.JSON(respond.Default.
				SetStatusText("Failed!").
				SetStatusCode(400).
				RespondWithMessage("Upload failed. Please try again later!"))
		}

		protoUser.Avatar = avatarName
	}

	if fullname != "" {
		protoUser.Fullname = fullname
	}

	mCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	response, err := grpc.UserServiceClient.UpdateUser(mCtx, &proto.UpdateUserRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Result: protoUser,
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.Error(500, 5445))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}

func UpdatePassword(ctx *gin.Context) {

	var (
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
		token        = ctx.GetHeader("Authorization")
		request      = &requests.UpdatePasswordRequest{
			Password:                ctx.PostForm("password"),
			NewPassword:             ctx.PostForm("new_password"),
			NewPasswordConfirmation: ctx.PostForm("new_password_confirmation"),
		}
	)
	defer cancel()

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	response, err := grpc.UserServiceClient.UpdatePassword(mCtx, &proto.UpdatePasswordRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		CurrentPassword:   request.Password,
		NewPassword:       request.NewPassword,
		VerifyNewPassword: request.NewPasswordConfirmation,
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusText("Failed").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Invalid Credentials!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusText("Success").
		SetStatusCode(http.StatusOK).
		RespondWithMessage("Password updated successfully!"))
	return
}
