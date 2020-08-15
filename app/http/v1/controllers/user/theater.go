package user

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/models"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"log"
	"net/http"
	"time"
)

// Get the current user from request
func Theater(ctx *gin.Context) {

	var (
		token = ctx.Request.Header.Get("Authorization")
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
	if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(code, result)
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}

func GetMediaSources(ctx *gin.Context) {

	var (
		mediaSources = make([]*proto.MediaSource, 0)
		token = ctx.Request.Header.Get("Authorization")
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
	)

	response, err := grpc.TheaterServiceClient.GetMediaSources(mCtx, &proto.MediaSourceAuthRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(code, result)
		return
	}

	if response.Result != nil {
		mediaSources = response.Result
	}

	ctx.JSON(respond.Default.SetStatusText("success").
		SetStatusCode(http.StatusOK).
		RespondWithResult(mediaSources))
	return
}

func SelectNewMediaSource(ctx *gin.Context)  {

	var (
		rules = govalidator.MapData{
			"source_id": []string{"required"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		token = ctx.Request.Header.Get("Authorization")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	_, err := grpc.TheaterServiceClient.SelectMediaSource(ctx, &proto.MediaSourceAuthRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Media: &proto.MediaSource{
			Id: ctx.PostForm("source_id"),
		},
	})

	if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(code, result)
		return
	}

	ctx.JSON(respond.Default.UpdateSucceeded())
	return

}

func ParseMediaSourceUri(ctx *gin.Context)  {

	var (
		rules = govalidator.MapData{
			"source": []string{"required", "media_source_uri"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	mediaSource := models.NewMediaSource(ctx.PostForm("source"))

	err := mediaSource.Parse()
	if err != nil {
		log.Println(err)
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
		token = ctx.Request.Header.Get("Authorization")
		rules = govalidator.MapData{
			"source": []string{"required", "media_source_uri"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
		validations := components.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	mediaSource := models.NewMediaSource(ctx.PostForm("source"))

	err := mediaSource.Parse()
	if err != nil {
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

		//index := ctx.PostForm("index")
		//if index == "" {
		//	ctx.JSON(respond.Default.Succeed(map[string] interface{} {
		//		"message": "Index parameter is required in torrent media source!",
		//		"code":    2342,
		//		"data":    mediaSource,
		//	}))
		//	return
		//}
		//
		//intI, err := strconv.Atoi(index)
		//if err != nil {
		//	ctx.JSON(respond.Default.SetStatusText("Failed").
		//		SetStatusCode(http.StatusBadRequest).
		//		RespondWithMessage("Index parameter is invalid!"))
		//	return
		//}
		//
		//var mediaFile *models.MediaFile
		//for index := range mediaSource.Files {
		//	if index == intI {
		//		mediaFile = &mediaSource.Files[intI]
		//	}
		//}
		//
		//if mediaFile == nil {
		//	ctx.JSON(respond.Default.SetStatusText("Failed").
		//		SetStatusCode(http.StatusBadRequest).
		//		RespondWithMessage("Index does not exists!"))
		//	return
		//}
		//
		//mediaFile.Download()
		//
		//log.Println(mediaFile)
		//log.Println("Downloading torrent from index: [", intI, "] ...")
		//return

	}

	protoMsg := mediaSource.Proto()

	if title := ctx.PostForm("title"); title != "" {
		protoMsg.Title = title
	}

	_, err = grpc.TheaterServiceClient.AddMediaSource(ctx, &proto.MediaSourceAuthRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Media: protoMsg,
	})

	if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
		ctx.JSON(code, result)
		return
	}

	ctx.JSON(respond.Default.Succeed(mediaSource.Proto()))
	return
}
