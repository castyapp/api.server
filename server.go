package main

import (
	"flag"
	"fmt"
	"github.com/CastyLab/api.server/app"
	"github.com/CastyLab/api.server/app/http/v1/middlewares"
	"github.com/CastyLab/api.server/app/http/v1/validators"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/thedevsaddam/govalidator"
	"log"
	"os"
	"time"
)

func main() {

	log.SetFlags(log.Lshortfile)

	port := flag.Int("port", 9002, "Casty http api port")
	flag.Parse()

	if err := sentry.Init(sentry.ClientOptions{ Dsn: os.Getenv("SENTRY_DSN") }); err != nil {
		log.Fatal(err)
	}

	defer sentry.Flush(5 * time.Second)

	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("APP_ENVIRONMENT") == "dev" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middlewares.CORSMiddleware)

	// register unique validator
	govalidator.AddCustomRule("access", validators.Access)
	govalidator.AddCustomRule("media_source_uri", validators.MediaSourceUri)

	app.RegisterRoutes(router)

	log.Printf("Server running and listening on port [%d]", *port)
	if err := router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}
