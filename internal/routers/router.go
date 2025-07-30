package routers

import (
	"Nuxus/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/register", controller.Register)
			user.POST("/verify-register", controller.VerifyRegister)
			user.GET("/login", controller.Login)
			user.POST("/password/reset", controller.RequestReset)
			user.POST("/password/verify-reset", controller.VerifyReset)
		}

	}

	return r
}
