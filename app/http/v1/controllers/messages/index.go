package messages

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"movie-night/grpc"
	"movie-night/proto"
	"net/http"
)

func Messages(ctx *gin.Context)  {

	var (
		receiverId = ctx.Param("receiver_id")
		token      = ctx.Request.Header.Get("Authorization")
	)

	response, err := grpc.MessagesServiceClient.GetUserMessages(ctx, &proto.GetMessagesRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		ReceiverId: receiverId,
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(500).
			RespondWithMessage("Could not get messages!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}
