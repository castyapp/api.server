package user

import (
	"context"
	"github.com/CastyLab/api.server/grpc"
	proto "github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Notifications(ctx *gin.Context)  {

	var (
		notifications = make([]*messages.Notification, 0)
		token = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
	)

	response, err := grpc.UserServiceClient.GetNotifications(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not get notifications!"))
		return
	}

	if response.Result != nil {
		notifications = response.Result
	}

	result := map[string] interface{} {
		"notifications": notifications,
		"unread_count": response.UnreadCount,
	}

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(result))
	return
}

func ReadAllNotifications(ctx *gin.Context)  {

	var (
		token = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
	)

	response, err := grpc.UserServiceClient.ReadAllNotifications(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil || response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not update notifications!"))
		return
	}

	ctx.JSON(respond.Default.UpdateSucceeded())
	return
}