package AirbnbSit

import (
	"AirbnbSit/authorization"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	authorization.SetupRoutes(r)
	r.Run(":8080")
}
