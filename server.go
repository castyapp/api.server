package main

import (
	"github.com/CastyLab/api.server/app"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

func main() {

	if err := sentry.Init(sentry.ClientOptions{ Dsn: os.Getenv("SENTRY_DSN") }); err != nil {
		log.Fatal(err)
	}

	defer sentry.Flush(5 * time.Second)

	// create a new application
	application := new(app.Application)

	// register providers
	application.RegisterProviders()

	// register router and serve
	application.RegisterAndServeRouter()
}
