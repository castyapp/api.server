package theater

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/MrJoshLab/go-respond"
	//"github.com/thedevsaddam/govalidator"
	//"github.com/CastyLab/api.server/app/components"
	//"github.com/CastyLab/api.server/app/components/strings"
	//"github.com/CastyLab/api.server/app/models"
	//"github.com/CastyLab/api.server/db"
	//"net/http"
)

func Subtitle(ctx *gin.Context)  {

	//var (
	//
	//	database = db.Connection
	//	user     = ctx.MustGet("user").(*models.User)
	//
	//	rules = govalidator.MapData{
	//		"lang":           []string{"required", "min:4", "max:30"},
	//		"file:subtitle":  []string{"required", "ext:srt", "size:20000000"},
	//		"theater_id":     []string{"required"},
	//	}
	//
	//	opts = govalidator.Options{
	//		Request:         ctx.Request,
	//		Rules:           rules,
	//		RequiredDefault: true,
	//	}
	//)
	//
	//if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {
	//
	//	validations := components.GetValidationErrorsFromGoValidator(validate)
	//	ctx.JSON(respond.Default.ValidationErrors(validations))
	//	return
	//}
	//
	//var (
	//	count int
	//	theater = new(models.Theater)
	//)
	//
	//database.Where(map[string] interface{} {
	//	"id": ctx.Request.Form.Get("theater_id"),
	//	"user_id": user.ID,
	//}).Find(theater).Count(&count)
	//
	//if count == 0 {
	//
	//	ctx.JSON(respond.Default.SetStatusCode(http.StatusNotFound).
	//		SetStatusText("Failed!").
	//		RespondWithMessage("Could not found the theater!"))
	//	return
	//}
	//
	//database.Model(theater).Related(&theater.Movie)
	//
	//var randomSubtitleName string
	//
	//if _, subtitleFile, err := ctx.Request.FormFile("subtitle"); err == nil {
	//
	//	randomSubtitleName = strings.RandomNumber(20)
	//
	//	if err := ctx.SaveUploadedFile(subtitleFile, fmt.Sprintf("./storage/uploads/subtitles/%s.srt", randomSubtitleName)); err != nil {
	//
	//		ctx.JSON(respond.Default.
	//			SetStatusText("Failed!").
	//			SetStatusCode(400).
	//			RespondWithMessage("Upload failed. Please try again later!"))
	//		return
	//	}
	//}
	//
	//subtitle := &models.Subtitle{
	//	Lang:     ctx.Request.Form.Get("lang"),
	//	File:     randomSubtitleName,
	//	Movie:    theater.Movie,
	//}
	//
	//if database = database.Create(subtitle); database.Error != nil {
	//
	//	ctx.JSON(respond.Default.Error(500, 5445))
	//	return
	//}
	//
	//ctx.JSON(respond.Default.Succeed(subtitle))
	//return
}
