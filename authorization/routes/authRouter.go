package routes

import (
	controller "auth/controllers"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.CORSMiddleware())
	incomingRoutes.POST("/signup", controller.Signup())
	incomingRoutes.POST("/login", controller.Login())
	incomingRoutes.POST("/verify-account", controller.VerifyAccount())
	incomingRoutes.POST("/password-code", controller.CreatePasswordRecoveryCode())
	incomingRoutes.POST("/forgot-password", controller.ForgotPassword())
	incomingRoutes.DELETE("/delete-account", controller.DeleteAccount())
}
