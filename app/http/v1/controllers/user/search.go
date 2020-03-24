package user

import (
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Search(ctx *gin.Context) {

	var (
		keyword = ctx.Query("keyword")
		rules = govalidator.MapData{
			"keyword":    []string{"required", "min:3", "max:20"},
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

	response, err := grpc.UserServiceClient.Search(ctx, &proto.SearchUserRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Keyword: keyword,
	})

	if err != nil || response.Code != http.StatusOK {

		ctx.JSON(respond.Default.NotFound())
		return
	}

	if response.Result == nil {
		ctx.JSON(respond.Default.Succeed([]interface{} {}))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}
