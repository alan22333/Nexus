package routers

import (
	"Nuxus/internal/controller"
	"Nuxus/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandler(), middleware.CORSMiddleware(), middleware.Recovery())

	v1 := r.Group("/api/v1")
	{
		// 普通路由
		user := v1.Group("/users")
		{
			user.POST("/register", controller.Register)
			user.POST("/verify-register", controller.VerifyRegister)
			user.POST("/login", controller.Login)
			user.POST("/password/reset", controller.RequestReset)
			user.POST("/password/verify-reset", controller.VerifyReset)
		}

		post := v1.Group("/posts")
		{
			post.GET("/", controller.ListPosts)
			post.GET("/popular", controller.ListPopularPosts)
			post.GET("/:id", controller.GetPost)

			comment := post.Group("/:id/comments")
			{
				comment.GET("/", controller.ListComment)
			}
		}

		tag := v1.Group("/tags")
		{
			tag.GET("/", controller.ListTags)
		}

		// 鉴权路由
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth())
		{
			me := auth.Group("me")
			{
				me.GET("/", controller.GetProfile)
				me.PUT("/", controller.UpdateProfile)
				me.POST("/avatar", controller.UpdateAvatar)

				// accout := auth.Group("/account")
				// {
				// 	accout.POST("/password", controller.UpdatePassword)
				// }
			}

			post := auth.Group("/posts")
			{
				post.POST("/", controller.CreatePost)
				post.PUT("/:id", controller.UpdatePost)
				post.DELETE("/:id", controller.DeletePost)
				post.GET("/:id/user-status",controller.GetUserStatus)

				comment := post.Group("/:id/comments")
				{
					comment.POST("/", controller.CreateComment)
				}

				like := post.Group("/:id/like")
				{
					like.POST("/", controller.LikePost)
				}
				favorite := post.Group("/:id/favorite")
				{
					favorite.POST("/", controller.FavoritePost)
				}
			}
			auth.DELETE("/comments/:commentId", controller.DeleteComment)

		}
	}

	return r
}
