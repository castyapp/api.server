package theater

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/subtitle"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"log"
	"net/http"
	"time"
)

// Get theater subtitles from grpc.theater service
func Subtitles(ctx *gin.Context) {

	var (
		subtitles  = make([]*proto.Subtitle, 0)
		mCtx, _    = context.WithTimeout(ctx, 10 * time.Second)
	)

	req := &proto.MediaSourceAuthRequest{
		Media: &proto.MediaSource{
			Id: ctx.Param("id"),
		},
	}

	if token := ctx.Request.Header.Get("Authorization"); token != "" {
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
	return
}

// Send request to grpc for adding subtitle to theater
func AddSubtitle(ctx *gin.Context) {

	var (
		rules = govalidator.MapData{
			"lang":           []string{"required", "min:4", "max:30"},
			"file:subtitle":  []string{"required", "ext:srt", "size:20000000"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		token      = ctx.Request.Header.Get("Authorization")
		mCtx, _    = context.WithTimeout(ctx, 10 * time.Second)
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	if subtitleFile, err := ctx.FormFile("subtitle"); err == nil {

		filename, err := subtitle.Save(subtitleFile)
		if err != nil {
			log.Println(err)
			sentry.CaptureException(err)
			ctx.JSON(respond.Default.
				SetStatusText("Failed!").
				SetStatusCode(400).
				RespondWithMessage("Upload failed. Please try again later!"))
			return
		}

		response, err := grpc.TheaterServiceClient.AddSubtitles(mCtx, &proto.AddSubtitlesRequest{
			MediaSourceId: ctx.Param("id"),
			Subtitles: []*proto.Subtitle{
				{
					Lang: ctx.PostForm("lang"),
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

		if response.Code != http.StatusOK {
			ctx.JSON(respond.Default.SetStatusText("failed").
				SetStatusCode(http.StatusBadRequest).
				RespondWithMessage("Could not add subtitle, please try again later!"))
			return
		}

		ctx.JSON(respond.Default.InsertSucceeded())
		return
	}

	ctx.JSON(respond.Default.SetStatusText("failed").
		SetStatusCode(http.StatusBadRequest).
		RespondWithMessage("Could not add subtitle, please try again later!"))
	return
}

// Send request to grpc for removing subtitle from theater
func RemoveSubtitle(ctx *gin.Context) {

	var (
		token      = ctx.Request.Header.Get("Authorization")
		mCtx, _    = context.WithTimeout(ctx, 10 * time.Second)
	)

	response, err := grpc.TheaterServiceClient.RemoveSubtitle(mCtx, &proto.RemoveSubtitleRequest{
		SubtitleId: ctx.Param("subtitle_id"),
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

	if response.Code != http.StatusOK {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	ctx.JSON(respond.Default.DeleteSucceeded())
	return
}