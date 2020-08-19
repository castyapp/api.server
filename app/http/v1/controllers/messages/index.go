package messages

import (
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
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

	if err != nil {
		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}
