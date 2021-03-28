package user

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/gin-gonic/gin"
)

func UpdateConnection(ctx *gin.Context) {

	var (
		connections  = make([]*proto.Connection, 0)
		token        = ctx.GetHeader("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
	)

	defer cancel()

	var service proto.Connection_Type
	switch serviceName := ctx.Param("service"); serviceName {
	case "google":
		service = proto.Connection_GOOGLE
	case "spotify":
		service = proto.Connection_SPOTIFY
	default:
		service = proto.Connection_UNKNOWN
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Invalid connection type!"))
		return
	}

	response, err := grpc.UserServiceClient.UpdateConnection(mCtx, &proto.ConnectionRequest{
		Connection: &proto.Connection{
			Type: service,
		},
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

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("failed").
			RespondWithMessage("Could not update connection."))
		return
	}

	if response.Result != nil {
		connections = response.Result
	}

	ctx.JSON(respond.Default.Succeed(connections))
	return

}

func GetConnection(ctx *gin.Context) {

	var (
		connections  = make([]*proto.Connection, 0)
		token        = ctx.GetHeader("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
	)

	defer cancel()

	var service proto.Connection_Type
	switch serviceName := ctx.Param("service"); serviceName {
	case "google":
		service = proto.Connection_GOOGLE
	case "spotify":
		service = proto.Connection_SPOTIFY
	default:
		service = proto.Connection_UNKNOWN
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Invalid connection type!"))
		return
	}

	response, err := grpc.UserServiceClient.GetConnection(mCtx, &proto.ConnectionRequest{
		Connection: &proto.Connection{
			Type: service,
		},
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

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("failed").
			RespondWithMessage("Could not get connection."))
		return
	}

	if response.Result != nil {
		connections = response.Result
	}

	ctx.JSON(respond.Default.Succeed(connections))
	return

}

func GetConnections(ctx *gin.Context) {

	var (
		connections  = make([]*proto.Connection, 0)
		token        = ctx.GetHeader("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
	)
	defer cancel()

	response, err := grpc.UserServiceClient.GetConnections(mCtx, &proto.AuthenticateRequest{
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
			RespondWithMessage("Could not get connections."))
		return
	}

	if response.Result != nil {
		connections = response.Result
	}

	ctx.JSON(respond.Default.Succeed(connections))
	return

}
