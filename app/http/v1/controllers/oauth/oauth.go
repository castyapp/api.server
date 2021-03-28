package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func Callback(ctx *gin.Context) {

	var (
		request = &requests.OauthCallbackRequest{
			Code: ctx.PostForm("code"),
		}
		service proto.Connection_Type
		token   = ctx.GetHeader("Authorization")
	)

	switch serviceName := ctx.Param("service"); serviceName {
	case "google":
		service = proto.Connection_GOOGLE
	case "spotify":
		service = proto.Connection_SPOTIFY
	default:
		service = proto.Connection_UNKNOWN
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Invalid OAUTH Service!"))
		return
	}

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	mCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &proto.OAUTHRequest{
		Code:    ctx.PostForm("code"),
		Service: service,
	}

	if token != "" {
		req.AuthRequest = &proto.AuthenticateRequest{
			Token: []byte(token),
		}
	}

	response, err := grpc.AuthServiceClient.CallbackOAUTH(mCtx, req)
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
	return
}
