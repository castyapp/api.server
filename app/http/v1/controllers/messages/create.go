package messages

import (
	"net/http"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/gin-gonic/gin"
)

func Create(ctx *gin.Context) {

	var (
		request = &requests.CreateMessageRequest{
			Content: ctx.PostForm("content"),
		}
		receiverId = ctx.Param("receiver_id")
		token      = ctx.GetHeader("Authorization")
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
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

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(response.Result))
	return
}
