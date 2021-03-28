package theater

import (
	"net/http"
	"strconv"

	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Update(ctx *gin.Context) {

	var (
		req = &requests.UpdateTheaterRequest{
			Description: ctx.PostForm("description"),
		}
		token = ctx.GetHeader("Authorization")
	)

	privacyInt, err := strconv.Atoi(ctx.PostForm("privacy"))
	if err == nil {
		req.Privacy = proto.PRIVACY(privacyInt)
	}

	videoPlayerAccessInt, err := strconv.Atoi(ctx.PostForm("video_player_access"))
	if err == nil {
		req.VideoPlayerAccess = proto.VIDEO_PLAYER_ACCESS(videoPlayerAccessInt)
	}

	if err := validator.New().Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	response, err := grpc.TheaterServiceClient.UpdateTheater(ctx, &proto.TheaterAuthRequest{
		Theater: &proto.Theater{
			Description:       req.Description,
			Privacy:           req.Privacy,
			VideoPlayerAccess: req.VideoPlayerAccess,
		},
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.UpdateFailed())
		return
	}

	ctx.JSON(respond.Default.UpdateSucceeded())
	return
}
