package user

import (
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

// Create a new user
func Create(c *gin.Context)  {

	var (
		rules = govalidator.MapData{
			"fullname":    []string{"min:4", "max:30"},
			"password":    []string{"required", "min:4", "max:30"},
			"username":    []string{"required", "between:3,8"},
			"email":       []string{"required", "min:4", "email"},
		}
		opts = govalidator.Options{
			Request:         c.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := components.GetValidationErrorsFromGoValidator(validate)
		c.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	response, err := grpc.UserServiceClient.CreateUser(c, &proto.CreateUserRequest{
		User: &messages.User{
			Fullname: c.PostForm("fullname"),
			Username: c.PostForm("username"),
			Email:    c.PostForm("email"),
			Password: c.PostForm("password"),
		},
	})

	if response != nil && response.Code == 420 {

		valErrs := make(map[string] interface{})
		for _, verr := range response.ValidationError {
			valErrs[verr.Field] = verr.Errors
		}

		c.JSON(respond.Default.ValidationErrors(valErrs))
		return
	}

	if err != nil || response == nil || response.Code != http.StatusOK {
		c.JSON(respond.Default.SetStatusCode(420).
			SetStatusText("failed").
			RespondWithMessage("Could not create user."))
		return
	}

	c.JSON(respond.Default.Succeed(map[string] interface{} {
		"token": string(response.Token),
		"refreshed_token": string(response.Token),
		"type": "bearer",
	}))
	return
}
