package theater

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"movie-night/grpc"
	"movie-night/proto/messages"
	"net/http"
)

// Create a new Theater
func Get(ctx *gin.Context)  {

	response, err := grpc.TheaterServiceClient.GetTheater(ctx, &messages.Theater{
		Id: ctx.Param("theater_id"),
	})

	if err != nil || response.Code != http.StatusOK {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}