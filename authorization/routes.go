package authorization

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST(`/register`, Register)
	}
}
