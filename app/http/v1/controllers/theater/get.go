package theater

import (
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
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