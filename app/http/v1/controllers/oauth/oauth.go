package oauth

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"time"
)

func Callback(ctx *gin.Context)  {

	var (
		service proto.OAUTHRequest_Service
		rules = govalidator.MapData{
			"code": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: false,
		}
		validate = govalidator.New(opts).Validate()
	)

	switch serviceName := ctx.Param("service"); serviceName {
	case "google": service = proto.OAUTHRequest_Google
	case "discord": service = proto.OAUTHRequest_Discord
	default: service = proto.OAUTHRequest_Invalid
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Invalid OAUTH Service!"))
		return
	}

	if validate.Encode() == "" {

		mCtx, cancel := context.WithTimeout(ctx, 10 * time.Second)
		defer cancel()

		response, err := grpc.AuthServiceClient.CallbackOAUTH(mCtx, &proto.OAUTHRequest{
			Code: ctx.PostForm("code"),
			Service: service,
		})

		if err != nil {
			ctx.JSON(respond.Default.SetStatusCode(http.StatusUnauthorized).
				SetStatusText("Failed!").
				RespondWithMessage("Unauthorized!"))
			return
		}

		if response.Code == http.StatusOK {
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