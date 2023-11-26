package routes

import (
	controller "auth/controllers"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.CORSMiddleware())
	incomingRoutes.GET("/", controller.GetUsers()).Use(middleware.Authenticate("USER"))
	incomingRoutes.GET("/:user_id", controller.GetUser()).Use(middleware.Authenticate("USER"))
}
