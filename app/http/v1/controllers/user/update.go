package user

import (
	"context"
	"fmt"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/strings"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"time"
)

func Update(ctx *gin.Context)  {

	var (
		protoUser = new(proto.User)
		token = ctx.Request.Header.Get("Authorization")
		rules = govalidator.MapData{
			"fullname":     []string{"min:4", "max:30"},
			"file:avatar":  []string{"ext:jpg,jpeg,png", "size:2000000"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		fullname = ctx.PostForm("fullname")
		avatar   = ctx.PostForm("avatar")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	if avatarFile, err := ctx.FormFile("avatar"); err == nil {
		avatar = strings.RandomNumber(20)
		if err := ctx.SaveUploadedFile(avatarFile, fmt.Sprintf("./storage/uploads/avatars/%s.png", avatar)); err != nil {
			sentry.CaptureException(err)
			ctx.JSON(respond.Default.
				SetStatusText("Failed!").
				SetStatusCode(400).
				RespondWithMessage("Upload failed. Please try again later!"))
			return
		}
		protoUser.Avatar = avatar
	}

	if fullname != "" {
		protoUser.Fullname = fullname
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)
	response, err := grpc.UserServiceClient.UpdateUser(mCtx, &proto.UpdateUserRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Result: protoUser,
	})

	if err != nil || response.Code != http.StatusOK {
		ctx.JSON(respond.Default.Error(500, 5445))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}