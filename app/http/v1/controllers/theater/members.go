package theater

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"gitlab.com/movienight1/grpc.proto"
	"gitlab.com/movienight1/grpc.proto/messages"
	"movie-night/grpc"
	"net/http"
)

// Create a new Theater
func GetMembers(ctx *gin.Context)  {

	var (
		members   = make([]*messages.User, 0)
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