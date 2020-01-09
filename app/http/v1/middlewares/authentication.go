package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/iamalirezaj/go-respond"
	"log"
	"movie-night/grpc"
	"movie-night/proto"
	"net/http"
	"strings"
	"time"
)

func Authentication(ctx *gin.Context)  {

	token := strings.ReplaceAll(ctx.GetHeader("Authorization"), "Bearer ", "")
	if token == "" {

		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(422).
			SetStatusText("Failed!").
			RespondWithMessage("Token is required!"))
		return
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	response, err := grpc.UserServiceClient.GetUser(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil || !response.Result.IsActive {

		log.Println(response, err)

		ctx.AbortWithStatusJSON(respond.Default.Error(http.StatusUnauthorized, 3012))
		return
	}

	ctx.Set("user", response.Result)
	ctx.Next()
}
