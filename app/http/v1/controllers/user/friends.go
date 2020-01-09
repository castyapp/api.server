package user

import (
	"context"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"gitlab.com/movienight1/grpc.proto"
	"movie-night/grpc"
	"net/http"
	"time"
)

func GetFriend(ctx *gin.Context)  {

	var (
		friendId = ctx.Param("friend_id")
		token = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
	)

	response, err := grpc.UserServiceClient.GetFriend(mCtx, &proto.FriendRequest{
		FriendId: friendId,
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not get friend!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}

func GetFriends(ctx *gin.Context)  {

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	response, err := grpc.UserServiceClient.GetFriends(mCtx, &proto.AuthenticateRequest{
		Token: []byte(ctx.Request.Header.Get("Authorization")),
	})

	if err != nil {

		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not get friends!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}
