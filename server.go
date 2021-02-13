package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/CastyLab/api.server/app"
	"github.com/CastyLab/api.server/app/http/v1/middlewares"
	"github.com/CastyLab/api.server/app/http/v1/validators"
	"github.com/CastyLab/api.server/config"
	"github.com/CastyLab/api.server/grpc"
	"github.com/CastyLab/api.server/storage"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

var (
	port *int
	host *string
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	port = flag.Int("port", 9002, "api server port")
	host = flag.String("host", "0.0.0.0", "api server host")
	configFileName := flag.String("config-file", "config.yml", "config.yaml file")

	flag.Parse()
	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	if err := config.Load(*configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not load config: %v", err))
	}

	if err := grpc.Configure(); err != nil {
		log.Fatal(fmt.Errorf("could not configure grpc.client: %v", err))
	}

	if err := storage.Configure(); err != nil {
		log.Fatal(fmt.Errorf("could not configure s3 bucket storage client: %v", err))
	}

	if err := sentry.Init(sentry.ClientOptions{Dsn: config.Map.Secrets.SentryDsn}); err != nil {
		log.Fatal(fmt.Errorf("could not initilize sentry: %v", err))
	}
}

func main() {

	// Since sentry emits events in the background we need to make sure
	// they are sent before we shut down
	defer sentry.Flush(5 * time.Second)

	gin.SetMode(gin.ReleaseMode)
	if config.Map.App.Env == "dev" || config.Map.App.Debug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middlewares.CORSMiddleware)

	// register unique validator
	govalidator.AddCustomRule("access", validators.Access)
	govalidator.AddCustomRule("media_source_uri", validators.MediaSourceUri)

	app.RegisterRoutes(router)

	log.Printf("Server running and listening on port [%s:%d]", *host, *port)
	if err := router.Run(fmt.Sprintf("%s:%d", *host, *port)); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}
