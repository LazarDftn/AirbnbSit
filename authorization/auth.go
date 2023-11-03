package authorization

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Username string
	Password string
	Name     string
	LastName string
	Email    string
	Address  string
}

var users = make(map[string]User)

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := users[user.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Korisnik već postoji"})
		return
	}

	users[user.Username] = user
	c.JSON(http.StatusCreated, gin.H{"message": "Korisnik je uspešno registrovan"})
}
