package user

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
)

func Notifications(ctx *gin.Context) {

	var (
		notifications = make([]*proto.Notification, 0)
		token         = ctx.GetHeader("Authorization")
		mCtx, cancel  = context.WithTimeout(ctx, 20*time.Second)
	)
	defer cancel()

	response, err := grpc.UserServiceClient.GetNotifications(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Result != nil {
		notifications = response.Result
	}

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(map[string]interface{}{
			"notifications": notifications,
			"unread_count":  response.UnreadCount,
		}))
	return
}

func ReadAllNotifications(ctx *gin.Context) {

	var (
		token        = ctx.GetHeader("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
	)
	defer cancel()

	response, err := grpc.UserServiceClient.ReadAllNotifications(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("failed").
			RespondWithMessage("Could not real all notifications."))
		return
	}

	ctx.JSON(respond.Default.UpdateSucceeded())
	return
}
