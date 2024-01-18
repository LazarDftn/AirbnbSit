package routes

import (
	"auth/controllers"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.CORSMiddleware())
	incomingRoutes.GET("/", controllers.GetUsers()).Use(middleware.Authenticate("USER"))
	incomingRoutes.GET("/:user_id", controllers.GetUser()).Use(middleware.Authenticate("USER"))
}
