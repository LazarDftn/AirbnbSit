package routes

import (
	controller "github.com/LazarDftn/AirbnbSit/authorization/controllers"
	"github.com/LazarDftn/AirbnbSit/authorization/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
