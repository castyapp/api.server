package auth

import (
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
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

	if err != nil || response.Code != http.StatusOK {

		log.Println(err, response)

		ctx.JSON(respond.Default.SetStatusCode(401).
			SetStatusText("failed").
			RespondWithMessage("Could not refresh token!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(map[string] interface{} {
		"token": string(response.Token),
		"refreshed_token": string(response.RefreshedToken),
		"type": "bearer",
	}))
}