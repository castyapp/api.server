package oauth

import (
	"context"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"time"
)

func Callback(ctx *gin.Context)  {

	var (
		service proto.Connection_Type
		rules = govalidator.MapData{
			"code": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: false,
		}
		token    = ctx.GetHeader("Authorization")
		validate = govalidator.New(opts).Validate()
	)

	switch serviceName := ctx.Param("service"); serviceName {
	case "google":  service = proto.Connection_GOOGLE
	case "discord": service = proto.Connection_DISCORD
	case "spotify": service = proto.Connection_SPOTIFY
	default: service = proto.Connection_UNKNOWN
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Invalid OAUTH Service!"))
		return
	}

	if validate.Encode() == "" {

		mCtx, cancel := context.WithTimeout(ctx, 10 * time.Second)
		defer cancel()

		req := &proto.OAUTHRequest{
			Code: ctx.PostForm("code"),
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

		if response.Code == http.StatusOK {

			if response.Message != "" {
				ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
					SetStatusText("success").
					RespondWithMessage(response.Message))
				return
			}

			ctx.JSON(respond.Default.Succeed(map[string] interface{} {
				"token": string(response.Token),
				"refreshed_token": string(response.RefreshedToken),
				"type": "bearer",
			}))
			return
		}

		ctx.JSON(respond.Default.SetStatusCode(http.StatusUnauthorized).
			SetStatusText("Failed!").
			RespondWithMessage("Unauthorized!"))
		return
	}

	validations := components.GetValidationErrorsFromGoValidator(validate)
	ctx.JSON(respond.Default.ValidationErrors(validations))
	return
}
