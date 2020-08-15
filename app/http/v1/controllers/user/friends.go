package user

import (
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPendingFriendRequests(ctx *gin.Context)  {

	friendRequests := make([]*proto.FriendRequest, 0)

	response, err := grpc.UserServiceClient.GetPendingFriendRequests(ctx, &proto.AuthenticateRequest{
		Token: []byte(ctx.Request.Header.Get("Authorization")),
	})

	if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(code, result)
		return
	}

	if response.Result != nil {
		friendRequests = response.Result
	}

	ctx.JSON(respond.Default.Succeed(friendRequests))
	return
}

func GetFriend(ctx *gin.Context)  {

	var (
		friendId = ctx.Param("friend_id")
		token = ctx.Request.Header.Get("Authorization")
	)

	response, err := grpc.UserServiceClient.GetFriend(ctx, &proto.FriendRequest{
		FriendId: friendId,
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(code, result)
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}

func GetFriendRequest(ctx *gin.Context)  {

	var (
		requestID = ctx.Param("friend_id")
		token = ctx.Request.Header.Get("Authorization")
	)

	response, err := grpc.UserServiceClient.GetFriendRequest(ctx, &proto.FriendRequest{
		RequestId: requestID,
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if _, _, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(http.StatusNotFound).
			RespondWithMessage("Could not find friend request!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(response))
	return
}

func GetFriends(ctx *gin.Context)  {

	friends := make([]*proto.User, 0)
	response, err := grpc.UserServiceClient.GetFriends(ctx, &proto.AuthenticateRequest{
		Token: []byte(ctx.Request.Header.Get("Authorization")),
	})

	if _, _, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not get friends!"))
		return
	}

	if response.Result != nil {
		friends = response.Result
	}

	ctx.JSON(respond.Default.Succeed(friends))
	return
}
