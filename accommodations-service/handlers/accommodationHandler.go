package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/repositories"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger
	repo   *repositories.AccommodationRepo
}

type SignedDetails struct {
	Username  string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func NewAccommodationsHandler(l *log.Logger, r *repositories.AccommodationRepo) *AccommodationsHandler {
	return &AccommodationsHandler{l, r}
}

func (a *AccommodationsHandler) PostAccommodation(c *gin.Context) {
	var accomm *domain.Accommodation
	fmt.Print(accomm)
	if err := c.ShouldBindJSON(&accomm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := a.repo.Insert(accomm)
	e := json.NewEncoder(c.Writer)
	e.Encode(res)
}

func (a *AccommodationsHandler) GetAllAccommodations(c *gin.Context) {

	accommodations, err := a.repo.GetAll()
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if accommodations == nil {
		return
	}

	err = accommodations.ToJSON(c.Writer)
	if err != nil {
		http.Error(c.Writer, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}

}

func (a *AccommodationsHandler) GetAccommById(c *gin.Context) {

	id := c.Param("accomm_id")

	accommodation, err := a.repo.GetAccommById(id)
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	accommodation.ToJSON(c.Writer)

}

func (a *AccommodationsHandler) MiddlewareAccommodationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		accommodation := &domain.Accommodation{}
		err := accommodation.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, accommodation)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *AccommodationsHandler) MiddlewareContentTypeSet() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
	}
}

func (a *AccommodationsHandler) CORSMiddleware() gin.HandlerFunc {
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

func (a *AccommodationsHandler) Authorize(role string) gin.HandlerFunc { // Check if user is authorized as Host or Guest depending on the 'role' parameter
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token") // get the token from the client
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not logged in!"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(clientToken) // check if token is invalid or expired
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		// If access is available only for hosts check if the user is a host
		if role == "HOST" && claims.User_type != "HOST" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have access to Host privileges!"})
			c.Abort()
			return
		}

		// If access is available only for guests check if the user is a guest
		if role == "GUEST" && claims.User_type != "GUEST" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have access to Guest privileges!"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		msg = err.Error()
		return
	}
	return claims, msg
}
