package app

import (
	"fmt"
	"github.com/CastyLab/api.server/app/http/v1/middlewares"
	"github.com/CastyLab/api.server/app/http/v1/validators"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
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

	gin.SetMode(gin.ReleaseMode)

	if env := getEnvBool("APP_DEBUG"); env {
		gin.SetMode(gin.DebugMode)
	}

	a.router = gin.New()
	a.router.Use(middlewares.CORSMiddleware)

	// register unique validator
	govalidator.AddCustomRule("access", validators.Access)
	govalidator.AddCustomRule("media_source_uri", validators.MediaSourceUri)
}

func (a *Application) RegisterAndServeRouter(port int)  {
	a.RegisterRoutes(); {
		log.Printf("Server running and listening on port [%d]", port)
		if err := a.router.Run(fmt.Sprintf(":%d", port)); err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
	}
}