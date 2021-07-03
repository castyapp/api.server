package middlewares

import (
	"github.com/castyapp/api.server/config"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", config.Map.HTTP.Rules.AccessControlAllowOrigin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	c.Writer.Header().Set(
		"Access-Control-Allow-Headers",
		"h-captcha-response, "+
			"Content-Type, "+
			"Content-Length, "+
			"Accept-Encoding, "+
			"X-CSRF-Token, "+
			"Authorization, "+
			"Service-Authorization, "+
			"accept, "+
			"origin, "+
			"Cache-Control, "+
			"X-Requested-With",
	)

	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
