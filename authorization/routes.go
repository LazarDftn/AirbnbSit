package authorization

import (
	"auth/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST(`/register`, controllers.Signup())
	}
}
