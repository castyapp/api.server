package user

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strconv"
	"time"
)

func UpdateState(ctx *gin.Context)  {

	var (
		rules = govalidator.MapData{
			"state":    []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	state, _ := strconv.Atoi(ctx.Request.Form.Get("state"))
	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	response, err := grpc.UserServiceClient.UpdateState(mCtx, &proto.UpdateStateRequest{
		State: messages.PERSONAL_STATE(state),
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(ctx.Request.Header.Get("Authorization")),
		},
	})

	if err != nil && response.Code == http.StatusOK {

		ctx.JSON(respond.Default.UpdateFailed())
		return
	}

	ctx.JSON(respond.Default.UpdateSucceeded())
	return
}
