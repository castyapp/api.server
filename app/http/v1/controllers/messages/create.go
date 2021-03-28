package messages

import (
	"net/http"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

func Create(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"content": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		receiverId = ctx.Param("receiver_id")
		token      = ctx.Request.Header.Get("Authorization")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	response, err := grpc.MessagesServiceClient.CreateMessage(ctx, &proto.MessageRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Message: &proto.Message{
			Reciever: &proto.User{Id: receiverId},
			Content:  ctx.PostForm("content"),
		},
	})

	if err != nil {
		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not send the message!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(response.Result))
	return
}
