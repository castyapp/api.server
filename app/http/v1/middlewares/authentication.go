package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/iamalirezaj/go-respond"
	"movie-night/grpc"
	"movie-night/proto"
	"net/http"
	"time"
)

func Authentication(c *gin.Context)  {

	var tokenString = c.Request.Header.Get("Authorization")

	if tokenString == "" {
		c.AbortWithStatusJSON(respond.Default.SetStatusCode(422).
			SetStatusText("Failed!").
			RespondWithMessage("Token is required!"))
		return
	}

	mCtx, _ := context.WithTimeout(c, 20 * time.Second)

	response, err := grpc.UserServiceClient.GetUser(mCtx, &proto.AuthenticateRequest{
		Token: []byte(tokenString),
	})

	if err != nil || !response.Result.IsActive {
		c.AbortWithStatusJSON(respond.Default.Error(http.StatusUnauthorized, 3012))
		return
	}

	c.Set("user", response.Result)
	c.Next()
}
