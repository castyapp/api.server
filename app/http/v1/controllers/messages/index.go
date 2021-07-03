package messages

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/gin-gonic/gin"
)

func Messages(ctx *gin.Context) {

	var (
		receiverID = ctx.Param("receiver_id")
		token      = ctx.GetHeader("Authorization")
	)

	response, err := grpc.MessagesServiceClient.GetUserMessages(ctx, &proto.GetMessagesRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		ReceiverId: receiverID,
	})

	if err != nil {
		code, result, ok := components.ParseGrpcErrorResponse(err)
		if !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
}
