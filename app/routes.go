package app

import (
	"github.com/castyapp/api.server/app/http/v1/controllers/auth"
	"github.com/castyapp/api.server/app/http/v1/controllers/messages"
	"github.com/castyapp/api.server/app/http/v1/controllers/oauth"
	"github.com/castyapp/api.server/app/http/v1/controllers/theater"
	"github.com/castyapp/api.server/app/http/v1/controllers/user"
	"github.com/castyapp/api.server/app/http/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	router.Static("/uploads", "./storage/uploads")

	// set v1 as namespace for all routes
	v1 := router.Group("v1")

	// oauth routes
	oauthGroup := v1.Group("oauth")
	oauthGroup.POST(":service/@callback", oauth.Callback)

	// authentication routes
	authGroup := v1.Group("auth")
	authGroup.POST("@create", auth.Create)
	authGroup.PUT("@create", auth.Refresh)

	// user routes
	authUserGroup := v1.Group("user")

	// use authentication middleware for user group routes
	authUserGroup.Use(middlewares.Authentication)

	authUserGroup.GET("@me", user.GetMe)
	authUserGroup.PUT("@me", user.Update)
	authUserGroup.PUT("@password", user.UpdatePassword)

	// Theater and media sources routes
	authUserGroup.GET("@theater", theater.Theater)
	authUserGroup.PUT("@theater", theater.Update)
	authUserGroup.POST("@media/select", theater.SelectNewMediaSource)
	authUserGroup.GET("@media", theater.GetMediaSources)
	authUserGroup.POST("@media", theater.AddNewMediaSource)
	authUserGroup.DELETE("@media", theater.DeleteMediaSource)
	authUserGroup.POST("@media/parse", theater.ParseMediaSourceURI)

	// notifications routes
	notifsGroup := authUserGroup.Group("@notifications")
	notifsGroup.GET("", user.Notifications)
	notifsGroup.PUT("", user.ReadAllNotifications)

	// theater routes
	theatersGroup := authUserGroup.Group("@theaters")
	theatersGroup.GET("", theater.GetFollowedTheaters)
	theatersGroup.POST(":id/invite", theater.Invite)
	theatersGroup.GET(":id/follow", theater.Follow)
	theatersGroup.GET(":id/unfollow", theater.Unfollow)
	theatersGroup.GET(":id/subtitles", theater.Subtitles)
	theatersGroup.POST(":id/subtitles", theater.AddSubtitle)
	theatersGroup.DELETE(":id/subtitles/:subtitle_id", theater.RemoveSubtitle)

	// friends routes
	friendsGroup := authUserGroup.Group("@friends")
	friendsGroup.GET("", user.GetFriends)
	friendsGroup.GET("pending", user.GetPendingFriendRequests)

	// friend routes
	friendGroup := authUserGroup.Group("@friend")
	friendGroup.GET(":friend_id", user.GetFriend)
	friendGroup.GET(":friend_id/request", user.SendFriendRequest)
	friendGroup.GET(":friend_id/request/get", user.GetFriendRequest)
	friendGroup.POST("accept", user.AcceptFriendRequest)

	// messages routes
	messagesGroup := authUserGroup.Group("@messages")
	messagesGroup.GET(":receiver_id", messages.Messages)
	messagesGroup.POST(":receiver_id", messages.Create)

	// connections routes
	// get user's oauth connections
	connectionsGroup := authUserGroup.Group("@connections")
	connectionsGroup.GET("", user.GetConnections)
	connectionsGroup.GET(":service", user.GetConnection)
	connectionsGroup.PUT(":service", user.UpdateConnection)

	// search for a spesefic user
	authUserGroup.GET("@search", user.Search)

	// user routes without authentication
	userGroup := v1.Group("user")
	userGroup.POST("@create", user.Create)
	userGroup.GET("@theater/:id", theater.Theater)
	userGroup.GET("@theater/:id/subtitles", theater.Subtitles)
}
