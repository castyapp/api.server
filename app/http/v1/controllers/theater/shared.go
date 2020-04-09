package theater

import (
	"context"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Get shared theater of user
func GetSharedTheaters(ctx *gin.Context) {

	var (
		theaters = make([]*proto.Theater, 0)
		token = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
	)

	response, err := grpc.TheaterServiceClient.GetUserSharedTheaters(mCtx, &proto.GetAllUserTheatersRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.SetStatusCode(int(response.Code)).
			SetStatusText(response.Status).
			RespondWithMessage(response.Message))
		return
	}

	if response.Result != nil {
		theaters = response.Result
	}

	ctx.JSON(respond.Default.Succeed(theaters))
	return

}