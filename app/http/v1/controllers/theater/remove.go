package theater

import (
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Create a new Theater
func Remove(ctx *gin.Context)  {

	token := ctx.Request.Header.Get("Authorization")

	response, err := grpc.TheaterServiceClient.RemoveTheater(ctx, &proto.TheaterAuthRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Theater: &proto.Theater{
			Id: ctx.Param("theater_id"),
		},
	})

	if err != nil || response.Code != http.StatusOK {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}