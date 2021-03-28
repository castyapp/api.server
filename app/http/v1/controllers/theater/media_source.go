package theater

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/app/models"
	"github.com/castyapp/api.server/grpc"
	"github.com/gin-gonic/gin"
)

func GetMediaSources(ctx *gin.Context) {

	var (
		mediaSources = make([]*proto.MediaSource, 0)
		token        = ctx.GetHeader("Authorization")
		mCtx, cancel = context.WithTimeout(ctx, 20*time.Second)
	)

	defer cancel()

	response, err := grpc.TheaterServiceClient.GetMediaSources(mCtx, &proto.MediaSourceAuthRequest{
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

	if response.Result != nil {
		mediaSources = response.Result
	}

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(mediaSources))
	return
}

func DeleteMediaSource(ctx *gin.Context) {

	var (
		token   = ctx.GetHeader("Authorization")
		request = &requests.MediaSourceRequest{
			SourceId: ctx.Query("source_id"),
		}
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	_, err := grpc.TheaterServiceClient.RemoveMediaSource(ctx, &proto.MediaSourceRemoveRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		MediaSourceId: request.SourceId,
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.UpdateSucceeded())
	return
}

func SelectNewMediaSource(ctx *gin.Context) {

	var (
		token   = ctx.GetHeader("Authorization")
		request = &requests.MediaSourceRequest{
			SourceId: ctx.PostForm("source_id"),
		}
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	response, err := grpc.TheaterServiceClient.SelectMediaSource(ctx, &proto.MediaSourceAuthRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Media: &proto.MediaSource{
			Id: request.SourceId,
		},
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(response.Result[0]))
	return

}

func ParseMediaSourceUri(ctx *gin.Context) {

	request := &requests.NewMediaSourceRequest{
		Source: ctx.PostForm("media_source_uri"),
	}
	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	accessToken := ctx.GetHeader("Service-Authorization")
	mediaSource := models.NewMediaSource(request.Source, accessToken)

	if err := mediaSource.Parse(); err != nil {
		ctx.JSON(respond.Default.
			SetStatusText("Failed!").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Could not parse data. Please try again later!"))
		return
	}

	if mediaSource.IsUnknown() {
		ctx.JSON(respond.Default.
			SetStatusText("Failed!").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Invalid media source type!"))
		return
	}

	if mediaSource.IsTorrent() {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusNotAcceptable).
			SetStatusText("Failed").
			RespondWithMessage("Torrent links are not available yet!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(mediaSource.Proto()))
	return
}

func AddNewMediaSource(ctx *gin.Context) {

	var (
		token   = ctx.GetHeader("Authorization")
		request = &requests.NewMediaSourceRequest{
			Source: ctx.PostForm("media_source_uri"),
		}
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	accessToken := ctx.GetHeader("Service-Authorization")
	mediaSource := models.NewMediaSource(request.Source, accessToken)

	if err := mediaSource.Parse(); err != nil {
		ctx.JSON(respond.Default.
			SetStatusText("Failed!").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Could not parse data. Please try again later!"))
		return
	}

	if mediaSource.IsUnknown() {
		ctx.JSON(respond.Default.
			SetStatusText("Failed!").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Invalid media source type!"))
		return
	}

	// check for index files inside of torrent
	if mediaSource.IsTorrent() {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusNotAcceptable).
			SetStatusText("Failed").
			RespondWithMessage("Torrent links are not available yet!"))
		return
	}

	protoMsg := mediaSource.Proto()

	if title := ctx.PostForm("title"); title != "" {
		protoMsg.Title = title
	}

	response, err := grpc.TheaterServiceClient.AddMediaSource(ctx, &proto.MediaSourceAuthRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Media: protoMsg,
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	ctx.JSON(respond.Default.Succeed(response.Result[0]))
	return
}
