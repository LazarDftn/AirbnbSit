package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"profile-service/domain"
	"profile-service/repositories"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	logger *log.Logger
	repo   *repositories.ProfileRepo
}

type SignedDetails struct {
	Username  string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func NewProfileHandler(l *log.Logger, r *repositories.ProfileRepo) *ProfileHandler {
	return &ProfileHandler{l, r}
}

func (p *ProfileHandler) Signup(c *gin.Context) {

	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultInsertionNumber, insertErr := p.repo.Insert(user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, resultInsertionNumber)
}

func (p *ProfileHandler) GetAllProfiles(c *gin.Context) {

	profiles, err := p.repo.GetAll()
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if profiles == nil {
		return
	}

	err = profiles.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
	c.JSON(http.StatusOK, "SUCCESS")

}

func (p *ProfileHandler) GetProfile(c *gin.Context) {

	email := c.Param("email")

	profile, err := p.repo.GetProfile(email)

	if profile == nil {
		return
	}

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	err = profile.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *ProfileHandler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
