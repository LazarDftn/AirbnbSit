package middleware

import (
	helper "auth/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		// after authentication if the user action requires a role, check it
		if role == "HOST" || role == "GUEST" {
			helper.CheckUserType(c, role)
		}

		c.Set("username", claims.Username)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}
