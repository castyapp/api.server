package user

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/gin-gonic/gin"
)

func GetPendingFriendRequests(ctx *gin.Context) {

	friendRequests := make([]*proto.FriendRequest, 0)

	response, err := grpc.UserServiceClient.GetPendingFriendRequests(ctx, &proto.AuthenticateRequest{
		Token: []byte(ctx.GetHeader("Authorization")),
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Result != nil {
		friendRequests = response.Result
	}

	ctx.JSON(respond.Default.Succeed(friendRequests))
}

func GetFriend(ctx *gin.Context) {

	var (
		friendID = ctx.Param("friend_id")
		token    = ctx.GetHeader("Authorization")
	)

	response, err := grpc.UserServiceClient.GetFriend(ctx, &proto.FriendRequest{
		FriendId: friendID,
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

	ctx.JSON(respond.Default.Succeed(response.Result))
}

func GetFriendRequest(ctx *gin.Context) {

	var (
		requestID = ctx.Param("friend_id")
		token     = ctx.GetHeader("Authorization")
	)

	response, err := grpc.UserServiceClient.GetFriendRequest(ctx, &proto.FriendRequest{
		RequestId: requestID,
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

	ctx.JSON(respond.Default.Succeed(response))
}

func GetFriends(ctx *gin.Context) {

	friends := make([]*proto.User, 0)
	response, err := grpc.UserServiceClient.GetFriends(ctx, &proto.AuthenticateRequest{
		Token: []byte(ctx.GetHeader("Authorization")),
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Result != nil {
		friends = response.Result
	}

	ctx.JSON(respond.Default.Succeed(friends))
}
