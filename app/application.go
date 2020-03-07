package app

import (
	"github.com/CastyLab/api.server/app/http/v1/controllers/auth"
	"github.com/CastyLab/api.server/app/http/v1/controllers/messages"
	"github.com/CastyLab/api.server/app/http/v1/controllers/theater"
	"github.com/CastyLab/api.server/app/http/v1/controllers/user"
	"github.com/CastyLab/api.server/app/http/v1/middlewares"
	"github.com/CastyLab/api.server/app/http/v1/validators"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"log"
	"os"
)

type Application struct {
	router *gin.Engine
}

func getEnvBool(key string) bool {
	switch os.Getenv(key) {
	case "true", "True", "1", "yes", "on":
		return true
	}
	return false
}

func (a *Application) RegisterProviders() {

	a.router = gin.Default()
	a.router.Use(middlewares.CORSMiddleware)

	gin.SetMode(gin.ReleaseMode)

	if env := getEnvBool("APP_DEBUG"); env != true {
		gin.SetMode(gin.DebugMode)
	}

	// register unique validator
	govalidator.AddCustomRule("access", validators.Access)
}

func (a *Application) RegisterAndServeRouter()  {

	a.router.Static("/uploads", "./storage/uploads")

	v1 := a.router.Group("v1"); {

		authGroup := v1.Group("auth"); {
			authGroup.POST("@create", auth.Create)
			authGroup.PUT("@create", auth.Refresh)
		}

		authUserGroup := v1.Group("user").Use(middlewares.Authentication); {

			authUserGroup.GET("@me", user.GetMe)
			authUserGroup.PUT("@state", user.UpdateState)

			authUserGroup.POST("@friends", user.SendFriendRequest)
			authUserGroup.GET("@friends", user.GetFriends)
			authUserGroup.GET("@friends/:friend_id", user.GetFriend)

			// theater routes
			authUserGroup.POST("@theaters", theater.Create)
			authUserGroup.GET("@theaters", theater.Index)
			authUserGroup.GET("@theaters/:theater_id", theater.Get)
			authUserGroup.GET("@theaters/:theater_id/members", theater.GetMembers)

			authUserGroup.GET("@messages/:receiver_id", messages.Messages)
			authUserGroup.POST("@messages/:receiver_id", messages.Create)
			authUserGroup.GET("@search", user.Search)
		}

		userGroup := v1.Group("user"); {
			userGroup.POST("@create", user.Create)
		}
	}

	if err := a.router.Run(":9002"); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}
}