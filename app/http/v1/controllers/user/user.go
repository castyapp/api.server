package user

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
)

// Get the current user from request
func GetMe(c *gin.Context) {

	if user, exists := c.Get("user"); exists {
		c.JSON(respond.Default.Succeed(user))
		return
	}

	c.JSON(respond.Default.SetStatusText("Failed!").
		SetStatusCode(500).
		RespondWithMessage("User does not exists in context!"))
}
