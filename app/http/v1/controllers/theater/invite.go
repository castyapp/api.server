package theater

import (
	"context"
	"encoding/json"
	"github.com/CastyLab/api.server/app/http/v1/requests"
	"github.com/CastyLab/api.server/grpc"
	proto "github.com/CastyLab/grpc.proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Invite(ctx *gin.Context)  {

	var (
		token   = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 10 * time.Second)
		request = new(requests.InviteToTheaterRequest)
	)

	rawJson, err := ctx.GetRawData()
	if err := json.Unmarshal(rawJson, request); err != nil {
		ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
			"friend_ids": []string {
				"Could not get parameters from raw json data!",
			},
		}))
		return
	}

	if len(request.FriendIDs) == 0 {
		ctx.JSON(respond.Default.ValidationErrors(map[string] interface{} {
			"friend_ids": []string {
				"You should at least pass 1 friend id!",
			},
		}))
		return
	}

	response, err := grpc.TheaterServiceClient.Invite(mCtx, &proto.InviteFriendsTheaterRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		TheaterId: ctx.Param("theater_id"),
		FriendIds: request.FriendIDs,
	})

	if err != nil || response.Code != http.StatusOK {

		log.Println(err, response)

		ctx.JSON(respond.Default.Error(500, 5445))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(200).
		SetStatusText("success").
		RespondWithMessage("Invitations sent successfully!"))
	return
}