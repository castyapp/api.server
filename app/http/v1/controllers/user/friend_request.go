package user

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"time"
)

func SendFriendRequest(ctx *gin.Context) {

	var (
		friendId = ctx.Param("friend_id")
		mCtx, _  = context.WithTimeout(ctx, 20 * time.Second)
		token    = ctx.Request.Header.Get("Authorization")
	)

	response, err := grpc.UserServiceClient.SendFriendRequest(mCtx, &proto.FriendRequest{
		FriendId: friendId,
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	switch response.Code {
	case 409:
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(409).
			RespondWithMessage("Friend request sent already!"))
		return
	case http.StatusOK:
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

	mCtx, cancel := context.WithTimeout(ctx, 20 * time.Second)
	defer cancel()

	response, err := grpc.UserServiceClient.AcceptFriendRequest(mCtx, &proto.FriendRequest{
		RequestId: ctx.PostForm("request_id"),
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(ctx.Request.Header.Get("Authorization")),
		},
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
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