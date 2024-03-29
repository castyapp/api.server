package theater

import (
	"context"
	"net/http"
	"time"

	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/components/subtitle"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// Get theater subtitles from grpc.theater service
func Subtitles(ctx *gin.Context) {

	var (
		subtitles        = make([]*proto.Subtitle, 0)
		mCtx, cancelFunc = context.WithTimeout(ctx, 10*time.Second)
	)
	defer cancelFunc()

	req := &proto.MediaSourceAuthRequest{
		Media: &proto.MediaSource{
			Id: ctx.Param("id"),
		},
	}

	if token := ctx.GetHeader("Authorization"); token != "" {
		req.AuthRequest = &proto.AuthenticateRequest{
			Token: []byte(token),
		}
	}

	response, err := grpc.TheaterServiceClient.GetSubtitles(mCtx, req)

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Result != nil {
		subtitles = response.Result
	}

	ctx.JSON(respond.Default.Succeed(subtitles))
}

// Send request to grpc for adding subtitle to theater
func AddSubtitle(ctx *gin.Context) {

	var (
		token            = ctx.GetHeader("Authorization")
		mCtx, cancelFunc = context.WithTimeout(ctx, 10*time.Second)
		request          = &requests.AddSubtitleRequest{
			Lang: ctx.PostForm("lang"),
		}
	)
	defer cancelFunc()

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	subtitleFile, err := ctx.FormFile("subtitle")
	if err != nil {
		ctx.JSON(respond.Default.SetStatusText("failed").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Subtitle file is required!"))
		return
	}

	filename, err := subtitle.Save(subtitleFile)
	if err != nil {
		sentry.CaptureException(err)
		ctx.JSON(respond.Default.
			SetStatusText("Failed!").
			SetStatusCode(400).
			RespondWithMessage("Upload failed. Please try again later!"))
		return
	}

	_, err = grpc.TheaterServiceClient.AddSubtitles(mCtx, &proto.AddSubtitlesRequest{
		MediaSourceId: ctx.Param("id"),
		Subtitles: []*proto.Subtitle{
			{
				Lang: request.Lang,
				File: filename,
			},
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

	ctx.JSON(respond.Default.InsertSucceeded())
}

// Send request to grpc for removing subtitle from theater
func RemoveSubtitle(ctx *gin.Context) {

	var (
		token            = ctx.GetHeader("Authorization")
		mCtx, cancelFunc = context.WithTimeout(ctx, 10*time.Second)
	)
	defer cancelFunc()

	_, err := grpc.TheaterServiceClient.RemoveSubtitle(mCtx, &proto.RemoveSubtitleRequest{
		SubtitleId:    ctx.Param("subtitle_id"),
		MediaSourceId: ctx.Param("id"),
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

	ctx.JSON(respond.Default.DeleteSucceeded())
}
