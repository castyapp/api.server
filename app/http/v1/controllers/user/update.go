package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/strings"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/api.server/storage"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/thedevsaddam/govalidator"
)

func Update(ctx *gin.Context) {

	var (
		protoUser = new(proto.User)
		token     = ctx.Request.Header.Get("Authorization")
		rules     = govalidator.MapData{
			"file:avatar": []string{"ext:jpg,jpeg,png", "size:2000000"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		fullname = ctx.PostForm("fullname")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

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
		token = ctx.Request.Header.Get("Authorization")
		rules = govalidator.MapData{
			"current_password":    []string{"required"},
			"new_password":        []string{"required"},
			"verify_new_password": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	mCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	response, err := grpc.UserServiceClient.UpdatePassword(mCtx, &proto.UpdatePasswordRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		CurrentPassword:   ctx.PostForm("current_password"),
		NewPassword:       ctx.PostForm("new_password"),
		VerifyNewPassword: ctx.PostForm("verify_new_password"),
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
