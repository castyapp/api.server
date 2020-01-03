package auth

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/fatih/structs"
	"movie-night/app/components/discord"
)

var discordClient = discord.NewDiscordOAuth()

func DiscordRedirectUser(c *gin.Context) {
	discordClient.Redirect(c)
}

func DiscordVerifyUser(c *gin.Context)  {
	oauth, err := discordClient.Verify(c)
	if err != nil {
		log.Println(err)
	}

	user, err := discordClient.GetUser(oauth.AccessToken)
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, structs.Map(user))
}
