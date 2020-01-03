package theater

import (
	"context"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"movie-night/grpc"
	"movie-night/proto"
	"movie-night/proto/messages"
	"net/http"
	"time"
)

func Index(ctx *gin.Context)  {

	var (
		theaters = make([]*messages.Theater, 0)
		token = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
	)

	response, err := grpc.TheaterServiceClient.GetUserTheaters(mCtx, &proto.GetAllUserTheatersRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.SetStatusCode(int(response.Code)).
			SetStatusText(response.Status).
			RespondWithMessage(response.Message))
		return
	}

	if response.Result != nil {
		theaters = response.Result
	}

	ctx.JSON(respond.Default.Succeed(theaters))
	return
}