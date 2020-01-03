package main

import (
	_ "github.com/joho/godotenv/autoload"
	"movie-night/app"
)

func main() {

	// create a new application
	application := new(app.Application)

	// register providers
	application.RegisterProviders()

	// register router and serve
	application.RegisterAndServeRouter()
}
