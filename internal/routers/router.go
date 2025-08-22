package routers

import (
	"Nuxus/internal/controller"
	"Nuxus/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userController    *controller.UserController
	postController    *controller.PostController
	tagController     *controller.TagController
	middlewareManager *middleware.MiddlewareManager
}

func NewRouter(
	userController *controller.UserController,
	postController *controller.PostController,
	tagController *controller.TagController,
	middlewareManager *middleware.MiddlewareManager,
) *Router {
	return &Router{
		userController:    userController,
		postController:    postController,
		tagController:     tagController,
		middlewareManager: middlewareManager,
	}
}

func (router *Router) SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(router.middlewareManager.ErrorHandler(),
		router.middlewareManager.CORSMiddleware(),
		router.middlewareManager.Recovery())

	v1 := r.Group("/api/v1")
	{
		// 普通路由
		user := v1.Group("/users")
		{
			user.POST("/register", router.userController.Register)
			user.POST("/verify-register", router.userController.VerifyRegister)
			user.POST("/login", router.userController.Login)
			user.POST("/password/reset", router.userController.RequestReset)
			user.POST("/password/verify-reset", router.userController.VerifyReset)
		}

		post := v1.Group("/posts")
		{
			post.GET("/", router.postController.ListPosts)
			post.GET("/popular", router.postController.ListPopularPosts)
			post.GET("/:id", router.postController.GetPost)

			comment := post.Group("/:id/comments")
			{
				comment.GET("/", router.postController.ListComment)
			}
		}

		tag := v1.Group("/tags")
		{
			tag.GET("/", router.tagController.ListTags)
		}

		// 鉴权路由
		auth := v1.Group("")
		auth.Use(router.middlewareManager.JWTAuth())
		{
			me := auth.Group("me")
			{
				me.GET("/", router.userController.GetProfile)
				me.PUT("/", router.userController.UpdateProfile)
				me.POST("/avatar", router.userController.UpdateAvatar)
			}

			post := auth.Group("/posts")
			{
				post.POST("/", router.postController.CreatePost)
				post.PUT("/:id", router.postController.UpdatePost)
				post.DELETE("/:id", router.postController.DeletePost)
				post.GET("/:id/user-status", router.postController.GetUserStatus)

				comment := post.Group("/:id/comments")
				{
					comment.POST("/", router.postController.CreateComment)
				}

				like := post.Group("/:id/like")
				{
					like.POST("/", router.postController.LikePost)
				}
				favorite := post.Group("/:id/favorite")
				{
					favorite.POST("/", router.postController.FavoritePost)
				}
			}
			auth.DELETE("/comments/:commentId", router.postController.DeleteComment)

		}
	}

	return r
}
