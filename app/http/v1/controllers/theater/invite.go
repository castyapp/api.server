package theater

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
)

func Invite(ctx *gin.Context) {

	var (
		token            = ctx.Request.Header.Get("Authorization")
		request          = new(requests.InviteToTheaterRequest)
		mCtx, cancelFunc = context.WithTimeout(ctx, 10*time.Second)
	)
	defer cancelFunc()

	rawJson, err := ctx.GetRawData()
	if err := json.Unmarshal(rawJson, request); err != nil {
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"friend_ids": []string{
				"Could not get parameters from raw json data!",
			},
		}))
		return
	}

	if len(request.FriendIDs) == 0 {
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"friend_ids": []string{
				"You should at least pass 1 friend id!",
			},
		}))
		return
	}

	response, err := grpc.TheaterServiceClient.Invite(mCtx, &proto.InviteFriendsTheaterRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		TheaterId: ctx.Param("id"),
		FriendIds: request.FriendIDs,
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("failed").
			RespondWithMessage("Could not send invitations, Please tray again later!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(200).
		SetStatusText("success").
		RespondWithMessage("Invited!"))
	return
}
