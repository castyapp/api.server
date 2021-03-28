package user

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/gin-gonic/gin"
)

func SendFriendRequest(ctx *gin.Context) {

	var (
		friendId     = ctx.Param("friend_id")
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
		token        = ctx.GetHeader("Authorization")
	)
	defer cancel()

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
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"friend_id": []string{
				"Could not sent a friend request to this user!",
			},
		}))
		return
	}
}

func AcceptFriendRequest(ctx *gin.Context) {

	var (
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
		token        = ctx.GetHeader("Authorization")
		request      = &requests.AcceptFriendRequest{
			RequestId: ctx.PostForm("request_id"),
		}
	)
	defer cancel()

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	response, err := grpc.UserServiceClient.AcceptFriendRequest(mCtx, &proto.FriendRequest{
		RequestId: request.RequestId,
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
	case http.StatusOK:
		ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
			SetStatusText("success").
			RespondWithMessage("Friend request accepted successfully!"))
		return
	default:
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"friend_id": []string{
				"Could not accept friend request, Pleae try again later!",
			},
		}))
		return
	}

}
