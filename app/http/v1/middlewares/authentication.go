package middlewares

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
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

	mCtx, cancel := context.WithTimeout(ctx, 20 * time.Second)
	defer cancel()

	response, err := grpc.UserServiceClient.GetUser(mCtx, &proto.AuthenticateRequest{
		Token: []byte(token),
	})

	if err != nil {
		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if !response.Result.IsActive {
		ctx.AbortWithStatusJSON(respond.Default.Error(http.StatusUnauthorized, 3012))
		return
	}

	ctx.Set("user", response.Result)
	ctx.Next()
}
