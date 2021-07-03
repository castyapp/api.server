package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func Callback(ctx *gin.Context) {

	var (
		request = &requests.OauthCallbackRequest{
			Code: ctx.PostForm("code"),
		}
		grpcRequest = new(proto.OAUTHRequest)
		token       = ctx.GetHeader("Authorization")
	)

	switch serviceName := ctx.Param("service"); serviceName {
	case "google":
		grpcRequest.Service = proto.Connection_GOOGLE
	case "spotify":
		grpcRequest.Service = proto.Connection_SPOTIFY
	default:
		grpcRequest.Service = proto.Connection_UNKNOWN
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Invalid OAUTH Service!"))
		return
	}

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	grpcRequest.Code = request.Code
	mCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if token != "" {
		grpcRequest.AuthRequest = &proto.AuthenticateRequest{
			Token: []byte(token),
		}
	}

	response, err := grpc.AuthServiceClient.CallbackOAUTH(mCtx, grpcRequest)
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(respond.Default.SetStatusCode(http.StatusUnauthorized).
			SetStatusText("Failed!").
			RespondWithMessage("Unauthorized!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":         string(response.Token),
		"refresh_token": string(response.RefreshedToken),
		"type":          "bearer",
	}))
}
