package theater

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"movie-night/grpc"
	"movie-night/proto"
	"net/http"
)

// Create a new Theater
func GetMembers(ctx *gin.Context)  {

	var (
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

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}