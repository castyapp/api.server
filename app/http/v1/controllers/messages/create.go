package messages

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gitlab.com/movienight1/grpc.proto"
	"movie-night/app/components"
	"movie-night/grpc"
	"net/http"
)

func Create(ctx *gin.Context)  {

	var (
		rules = govalidator.MapData{
			"content": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		receiverId  = ctx.Param("receiver_id")
		token       = ctx.Request.Header.Get("Authorization")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	response, err := grpc.MessagesServiceClient.CreateMessage(ctx, &proto.CreateMessageRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		RecieverId: receiverId,
		Content: ctx.PostForm("content"),
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not send the message!"))
		return
	}

	// Send message through websocket

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(response.Result))
	return
}
