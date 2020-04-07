package user

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/api.server/internal"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
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
		friendId = ctx.PostForm("friend_id")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	response, err := grpc.UserServiceClient.SendFriendRequest(mCtx, &proto.FriendRequest{
		FriendId: friendId,
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(ctx.Request.Header.Get("Authorization")),
		},
	})

	if err != nil {
		ctx.JSON(respond.Default.SetStatusCode(500).
			SetStatusText("Failed!").
			RespondWithMessage("Something went wrong, Please try again later!"))
		return
	}

	switch response.Code {
	case 409:
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(409).
			RespondWithMessage("Friend request sent already!"))
		return
	case http.StatusOK:

		// send a new notification event to friend
		_ = internal.Client.UserService.SendNewNotificationsEvent(friendId)

		ctx.JSON(respond.Default.InsertSucceeded())
		return
	default:
		ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
			"friend_id": []string {
				"Could not sent a friend request to this user!",
			},
		}))
		return
	}
}

func AcceptFriendRequest(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"request_id": []string{"required"},
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

	response, err := grpc.UserServiceClient.AcceptFriendRequest(mCtx, &proto.FriendRequest{
		RequestId: ctx.PostForm("request_id"),
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(ctx.Request.Header.Get("Authorization")),
		},
	})

	if err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Something went wrong, Please try again later!"))
		return
	}

	switch response.Code {
	case http.StatusOK:
		ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
			SetStatusText("success").
			RespondWithMessage("Friend request accepted successfully!"))
		return
	default:
		ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
			"friend_id": []string {
				"Could not accept friend request, Pleae try again later!",
			},
		}))
		return
	}

}