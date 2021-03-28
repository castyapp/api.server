package user

import (
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/castyapp/api.server/app/components"
	"github.com/castyapp/api.server/app/http/v1/requests"
	"github.com/castyapp/api.server/app/http/v1/validators"
	"github.com/castyapp/api.server/grpc"
	"github.com/gin-gonic/gin"
)

func Search(ctx *gin.Context) {

	var (
		request = &requests.SearchUserRequest{
			Keyword: ctx.Query("keyword"),
		}
		token = ctx.GetHeader("Authorization")
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	response, err := grpc.UserServiceClient.Search(ctx, &proto.SearchUserRequest{
		AuthRequest: &proto.AuthenticateRequest{
			Token: []byte(token),
		},
		Keyword: request.Keyword,
	})

	if err != nil {
		if code, result, ok := components.ParseGrpcErrorResponse(err); !ok {
			ctx.JSON(code, result)
			return
		}
	}

	if response.Result == nil {
		ctx.JSON(respond.Default.Succeed([]interface{}{}))
		return
	}

	ctx.JSON(respond.Default.Succeed(response.Result))
	return
}
