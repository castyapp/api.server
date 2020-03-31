package theater

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/strings"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// Create a new Theater
func Create(ctx *gin.Context) {

	var (
		token = ctx.Request.Header.Get("Authorization")
		rules = govalidator.MapData{
			"title":               []string{"required", "min:4", "max:30"},
			"video_player_access": []string{"required", "bool"},
			"privacy":             []string{"required", "access"},
			"movie_uri":           []string{"required"},
			"file:poster":         []string{"ext:jpg,jpeg,png", "size:2000000"},
			"file:subtitles":      []string{"ext:srt", "size:20000000"},
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

	var moviePosterName = "default"
	if bannerFile, err := ctx.FormFile("poster"); err == nil {
		moviePosterName = strings.RandomNumber(20)
		if err := ctx.SaveUploadedFile(bannerFile, fmt.Sprintf("./storage/uploads/posters/%s.png", moviePosterName)); err != nil {

			sentry.CaptureException(err)

			ctx.JSON(respond.Default.
				SetStatusText("Failed!").
				SetStatusCode(400).
				RespondWithMessage("Upload failed. Please try again later!"))
			return
		}
	}

	privacy, _ := strconv.Atoi(ctx.PostForm("privacy"))
	videoPlayerAccess, _ := strconv.Atoi(ctx.PostForm("video_player_access"))

	movieType := proto.MovieType_UNKNOWN
	typeID, err := strconv.Atoi(ctx.PostForm("type"))
	if err == nil {
		movieType = proto.MovieType(typeID)
	}

	mCtx, _ := context.WithTimeout(ctx, 20*time.Second)
	response, err := grpc.TheaterServiceClient.CreateTheater(mCtx, &proto.CreateTheaterRequest{
		Theater: &proto.Theater{
			Title:             ctx.PostForm("title"),
			Privacy:           proto.PRIVACY(privacy),
			VideoPlayerAccess: proto.PRIVACY(videoPlayerAccess),
			Movie: &proto.Movie{
				Poster: moviePosterName,
				Type:   movieType,
				Uri:    ctx.PostForm("uri"),
			},
		},
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.Error(500, 5445))
		return
	}

	ctx.JSON(respond.Default.InsertSucceeded())
	return
}
