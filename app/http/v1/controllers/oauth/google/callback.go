package google

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	proto "github.com/CastyLab/grpc.proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"log"
	"net/http"
	"time"
)

func Callback(ctx *gin.Context)  {

	var (
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

	if validate.Encode() == "" {

		mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)
		response, err := grpc.AuthServiceClient.CallbackOAUTH(mCtx, &proto.OAUTHRequest{
			Code: ctx.PostForm("code"),
			Service: proto.OAUTHRequest_Google,
		})

		if err != nil {
			log.Println(err)
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