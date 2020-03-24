package theater

import (
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Create a new Theater
func GetMembers(ctx *gin.Context)  {

	var (
		members   = make([]*proto.User, 0)
		theaterId = ctx.Param("theater_id")
		tokne     = ctx.Request.Header.Get("Authorization")
	)

	response, err := grpc.TheaterServiceClient.GetMembers(ctx, &proto.GetTheaterMembersRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(tokne),
		},
		TheaterId: theaterId,
	})

	if err != nil || response.Code != http.StatusOK {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	if response.Result != nil {
		members = response.Result
	}

	ctx.JSON(respond.Default.Succeed(members))
	return
}