package auth

import (
	"github.com/CastyLab/api.server/app/components/oauth/discord"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
