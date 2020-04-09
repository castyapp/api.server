package app

import (
	"github.com/CastyLab/api.server/app/http/v1/controllers/auth"
	"github.com/CastyLab/api.server/app/http/v1/controllers/messages"
	"github.com/CastyLab/api.server/app/http/v1/controllers/oauth"
	"github.com/CastyLab/api.server/app/http/v1/controllers/theater"
	"github.com/CastyLab/api.server/app/http/v1/controllers/user"
	"github.com/CastyLab/api.server/app/http/v1/middlewares"
)

func (a *Application) RegisterRoutes()  {

	a.router.Static("/uploads", "./storage/uploads")

	v1 := a.router.Group("v1"); {

		oauthGroup := v1.Group("oauth"); {
			oauthGroup.POST(":service/@callback", oauth.Callback)
		}

		authGroup := v1.Group("auth"); {
			authGroup.POST("@create", auth.Create)
			authGroup.PUT("@create", auth.Refresh)
		}

		authUserGroup := v1.Group("user").Use(middlewares.Authentication); {

			authUserGroup.GET("@me", user.GetMe)
			authUserGroup.PUT("@me", user.Update)

			authUserGroup.GET("@notifications", user.Notifications)
			authUserGroup.DELETE("@notifications", user.ReadAllNotifications)

			authUserGroup.POST("@friends", user.SendFriendRequest)
			authUserGroup.GET("@friends", user.GetFriends)
			authUserGroup.GET("@friends/:friend_id", user.GetFriend)
			authUserGroup.GET("@friends/:friend_id/@fr", user.GetFriendRequest)
			authUserGroup.POST("@friends/accept", user.AcceptFriendRequest)

			// theater routes
			authUserGroup.POST("@theaters", theater.Create)
			authUserGroup.GET("@theaters", theater.Index)
			authUserGroup.GET("@shared_theaters", theater.GetSharedTheaters)
			authUserGroup.GET("@theaters/:theater_id", theater.Get)
			authUserGroup.DELETE("@theaters/:theater_id", theater.Remove)

			authUserGroup.GET("@theaters/:theater_id/subtitles", theater.Subtitles)
			authUserGroup.POST("@theaters/:theater_id/subtitles", theater.AddSubtitle)
			authUserGroup.DELETE("@theaters/:theater_id/subtitles/:subtitle_id", theater.RemoveSubtitle)

			authUserGroup.POST("@theaters/:theater_id/invite", theater.Invite)
			
			authUserGroup.GET("@messages/:receiver_id", messages.Messages)
			authUserGroup.POST("@messages/:receiver_id", messages.Create)
			authUserGroup.GET("@search", user.Search)
		}

		userGroup := v1.Group("user"); {
			userGroup.POST("@create", user.Create)
		}
	}
}