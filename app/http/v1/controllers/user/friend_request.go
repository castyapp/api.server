package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/iamalirezaj/go-respond"
	"github.com/thedevsaddam/govalidator"
	"gitlab.com/movienight1/grpc.proto"
	"movie-night/app/components"
	"movie-night/grpc"
	"net/http"
	"time"
)

func SendFriendRequest(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"friend_id": []string{"required"},
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

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	response, err := grpc.UserServiceClient.SendFriendRequest(mCtx, &proto.FriendRequest{
		FriendId: ctx.PostForm("friend_id"),
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(ctx.Request.Header.Get("Authorization")),
		},
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
			"friend_id": []string {
				"Could not sent a friend request to this user!",
			},
		}))
		return
	}

	ctx.JSON(respond.Default.InsertSucceeded())
	return
}
