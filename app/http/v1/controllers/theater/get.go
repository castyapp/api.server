package theater

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/gin-gonic/gin"
)

// Get the current user from request
func Theater(ctx *gin.Context) {

	var (
		token = ctx.GetHeader("Authorization")
		req   = new(proto.GetTheaterRequest)
	)

	if token != "" {
		req = &proto.GetTheaterRequest{
			AuthRequest: &proto.AuthenticateRequest{
				Token: []byte(token),
			},
		}
	}

	if user := ctx.Param("id"); user != "" {
		req.User = user
	}

	response, err := grpc.TheaterServiceClient.GetTheater(ctx, req)

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
}
