package auth

import (
	"context"
	"github.com/CastyLab/api.server/app/components"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/grpc.proto"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"log"
	"net/http"
	"time"
)

// Create user jwt token
func Create(c *gin.Context) {

	var (
		username string
		rules = govalidator.MapData{
			"password":    []string{"required", "min:4", "max:30"},
			"username":    []string{"between:3,8"},
			"email":       []string{"min:4", "email"},
		}
		opts = govalidator.Options{
			Request:         c.Request,
			Rules:           rules,
			RequiredDefault: false,
		}
		validator = govalidator.New(opts)
		validate  = validator.Validate()
	)

	if validate.Encode() == "" {

		if postFormUser := c.PostForm("username"); postFormUser != "" {
			username = postFormUser
		}

		if postFormUser := c.PostForm("email"); postFormUser != "" {
			username = postFormUser
		}

		if username == "" {
			c.JSON(respond.Default.ValidationErrors(map[string] interface{} {
				"user": []string {
					"Username or email is required!",
				},
			}))
			return
		}

		mCtx, _ := context.WithTimeout(c, 20 * time.Second)
		response, err := grpc.AuthServiceClient.Authenticate(mCtx, &proto.AuthRequest{
			User: username,
			Pass: c.PostForm("password"),
		})

		if err != nil {
			log.Println(err)
			return
		}

		if response.Code == http.StatusOK {
			c.JSON(respond.Default.Succeed(map[string] interface{} {
				"token": string(response.Token),
				"refreshed_token": string(response.RefreshedToken),
				"type": "bearer",
			}))
			return
		}

		c.JSON(respond.Default.SetStatusCode(http.StatusUnauthorized).
			SetStatusText("Failed!").
			RespondWithMessage("Unauthorized!"))
		return
	}

	validations := components.GetValidationErrorsFromGoValidator(validate)
	c.JSON(respond.Default.ValidationErrors(validations))
	return
}