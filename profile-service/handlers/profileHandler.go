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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	res := p.repo.CheckUsernameExists(*user.Username)

	if res == "error" {
		c.JSON(419, gin.H{"error": "Error searching for identical usernames!"})
		return
	}

	if res != "" {
		c.JSON(418, gin.H{"error": "This username already exists!"})
		return
	}

	result, insertErr := p.repo.Insert(user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
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

	id := c.Param("id")

	profile, err := p.repo.GetProfile(id)

	if profile == nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (p *ProfileHandler) DeleteProfile(c *gin.Context) {

	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	err = p.repo.Delete(objectId)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, "")
}

func (p *ProfileHandler) EditProfile(c *gin.Context) {

	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)

	var profile domain.User

	if err := c.BindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	if p.repo.CheckUsernameExists(*profile.Username) != "" {
		c.JSON(428, "This username already exists")
		return
	}

	err = p.repo.EditProfile(profile, objectId)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, "")
}

func (p *ProfileHandler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
