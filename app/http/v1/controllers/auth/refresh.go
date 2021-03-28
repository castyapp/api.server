package auth

import (
	"net/http"
	"strings"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
)

// Refresh expired authentication token
func Refresh(ctx *gin.Context) {

	token := strings.ReplaceAll(ctx.GetHeader("Authorization"), "Bearer ", "")
	if token == "" {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(422).
			SetStatusText("Failed!").
			RespondWithMessage("Token is required!"))
		return
	}

	response, err := grpc.AuthServiceClient.RefreshToken(ctx, &proto.RefreshTokenRequest{
		RefreshedToken: []byte(token),
	})

	if err != nil {
		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(401).
			SetStatusText("failed").
			RespondWithMessage("Could not refresh token!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":           string(response.Token),
		"refreshed_token": string(response.RefreshedToken),
		"type":            "bearer",
	}))
	return
}
