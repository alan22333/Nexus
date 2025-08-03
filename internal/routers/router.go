package routers

import (
	"Nuxus/internal/controller"
	"Nuxus/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// 普通路由
		user := v1.Group("/users")
		{
			user.POST("/register", controller.Register)
			user.POST("/verify-register", controller.VerifyRegister)
			user.GET("/login", controller.Login)
			user.POST("/password/reset", controller.RequestReset)
			user.POST("/password/verify-reset", controller.VerifyReset)
		}

		post := v1.Group("/posts")
		{
			post.GET("/", controller.ListPosts)
			post.GET("/popular", controller.ListPopularPosts)
			post.GET("/:id", controller.GetPost)
		}

		tag := v1.Group("/tags")
		{
			tag.GET("/", controller.ListTags)
		}

		// 鉴权路由
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth())
		{
			post := auth.Group("/posts")
			{
				post.POST("/", controller.CreatePost)
				post.PUT("/:id", controller.UpdatePost)
				post.DELETE("/:id", controller.DeletePost)
			}

		}
	}

	return r
}
