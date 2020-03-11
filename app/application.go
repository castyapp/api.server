package app

import (
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
	a.RegisterRoutes(); {
		if err := a.router.Run(":9002"); err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
	}
}