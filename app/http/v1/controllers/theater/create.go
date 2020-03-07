package theater

import (
	"context"
	"fmt"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/app/components/strings"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strconv"
	"time"
)

// Create a new Theater
func Create(ctx *gin.Context)  {

	var (
		token = ctx.Request.Header.Get("Authorization")
		rules = govalidator.MapData{
			"title":                 []string{"required", "min:4", "max:30"},
			"video_player_access":   []string{"required", "bool"},
			"privacy":               []string{"required", "access"},
			"movie_uri":             []string{"required", "url"},
			"file:poster":           []string{"ext:jpg,jpeg,png", "size:2000000"},
			"file:subtitles":        []string{"ext:srt", "size:20000000"},
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

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)
	response, err := grpc.TheaterServiceClient.CreateTheater(mCtx, &proto.CreateTheaterRequest{
		Theater: &messages.Theater{
			Title: ctx.PostForm("title"),
			Privacy: messages.PRIVACY(privacy),
			VideoPlayerAccess: messages.PRIVACY(videoPlayerAccess),
			Movie: &messages.Movie{
				Poster: moviePosterName,
				MovieUri: ctx.PostForm("movie_uri"),
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