package theater

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
)

func GetFollowedTheaters(ctx *gin.Context) {

	var (
		theaters     = make([]*proto.Theater, 0)
		token        = ctx.Request.Header.Get("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	)

	defer cancel()

	response, err := grpc.TheaterServiceClient.GetFollowedTheaters(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Result != nil {
		theaters = response.Result
	}

	ctx.JSON(respond.Default.Succeed(theaters))
	return
}

func Follow(ctx *gin.Context) {

	var (
		token        = ctx.Request.Header.Get("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	)

	defer cancel()

	response, err := grpc.TheaterServiceClient.Follow(mCtx, &proto.TheaterAuthRequest{
		Theater: &proto.Theater{
			Id: ctx.Param("id"),
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

	if response.Code == http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
			SetStatusText("Success").
			RespondWithMessage("Followed successfully!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
		SetStatusText("Success").
		RespondWithMessage("Could not follow. please try again later!"))
	return
}

func Unfollow(ctx *gin.Context) {

	var (
		token        = ctx.Request.Header.Get("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	)

	defer cancel()

	response, err := grpc.TheaterServiceClient.Unfollow(mCtx, &proto.TheaterAuthRequest{
		Theater: &proto.Theater{
			Id: ctx.Param("id"),
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

	if response.Code == http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
			SetStatusText("Success").
			RespondWithMessage("Unfollowed successfully!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
		SetStatusText("Success").
		RespondWithMessage("Could not unfollow. please try again later!"))
	return
}
