package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"reservation-service/domain"
	"reservation-service/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type KeyProduct struct{}

type ReservationsHandler struct {
	logger *log.Logger
	// NoSQL: injecting reservation repository
	repo *repositories.ReservationRepo
}

type SignedDetails struct {
	Username  string
	User_type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// Injecting the logger makes this code much more testable.
func NewReservationsHandler(l *log.Logger, r *repositories.ReservationRepo) *ReservationsHandler {
	return &ReservationsHandler{l, r}
}

func (r *ReservationsHandler) PostReservation(c *gin.Context) {
	var res *domain.Reservation
	fmt.Print(res)
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := r.repo.InsertReservation(res)
	if err != "" {
		r.logger.Print("Database exception: ", err)
		e := json.NewEncoder(c.Writer)
		e.Encode(err)
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

func (rr *ReservationsHandler) GetResByAccommodation(c *gin.Context) {

	location := c.Param("location")
	accommId := c.Param("accomm_id")

	reservations, err := rr.repo.GetReservationsByAccomm(location, accommId)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
	}

	if reservations == nil {
		return
	}

	err = reservations.ToJSON(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		rr.logger.Fatal("Unable to convert to json :", err)
		return
	}

}

func (r *ReservationsHandler) PostPrice(c *gin.Context) {
	var price *domain.AccommPrice
	fmt.Print(price)
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := r.repo.InsertAccommodationPrice(price)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

func (rr *ReservationsHandler) GetPriceByAccommodation(c *gin.Context) {

	accommId := c.Param("accomm_id")

	price, err := rr.repo.GetPriceByAccomm(accommId)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
	}

	if price == nil {
		return
	}

	c.JSON(http.StatusOK, price)

}

func (rr *ReservationsHandler) CheckPrice(c *gin.Context) {

	var res domain.Reservation
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variations := rr.repo.CheckPrice(res)
	if variations != nil {
		price, days := CompareAndCalculate(variations, res)
		retVal := struct {
			Price       int
			Percentages []domain.PriceVariation
			OverlapDays []int
		}{
			Price:       int(price),
			Percentages: variations,
			OverlapDays: days,
		}
		e := json.NewEncoder(c.Writer)
		e.Encode(retVal)
	} else {
		retVal := struct {
			Days int
		}{
			Days: int(GetNumberOfDays(res.StartDate, res.EndDate)),
		}
		e := json.NewEncoder(c.Writer)
		e.Encode(retVal)
	}

}

func (rr *ReservationsHandler) CreatePriceVariation(c *gin.Context) {

	var vr domain.PriceVariation

	if err := c.ShouldBindJSON(&vr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := rr.repo.InsertPriceVariation(&vr)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
		e := json.NewEncoder(c.Writer)
		e.Encode(message)
	}
	e := json.NewEncoder(c.Writer)
	e.Encode(message)
}

func (rr *ReservationsHandler) GetPriceVariationByAccommId(c *gin.Context) {
	location := c.Param("location")
	accommId := c.Param("accomm_id")

	variations, err := rr.repo.GetVariationsByAccommId(location, accommId)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
	}

	if variations == nil {
		return
	}

	e := json.NewEncoder(c.Writer)
	e.Encode(variations)
}

func (rr *ReservationsHandler) SearchAccommodations(c *gin.Context) {

	var av domain.Availability

	if err := c.ShouldBindJSON(&av); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*availabilites, err := rr.repo.SearchAccommodations(av)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
	}

	if availabilites == nil {
		return
	}

	e := json.NewEncoder(c.Writer)
	e.Encode(availabilites)*/
}

func CompareAndCalculate(vr []domain.PriceVariation, res domain.Reservation) (float64, []int) {

	var finalPrice = 0.0
	var daysToAddVariations []int
	var sumDaysToAddVariation float64

	resDays := GetNumberOfDays(res.StartDate, res.EndDate)

	for i := 0; i < len(vr); i++ {

		daysToAddVariation := math.Min(float64(res.EndDate.UnixMilli()/86400000), float64(vr[i].EndDate.UnixMilli()/86400000)) -
			math.Max(float64(res.StartDate.UnixMilli()/86400000), float64(vr[i].StartDate.UnixMilli()/86400000)) + 1

		percentageIncrease := (float64(100+vr[i].Percentage) / 100.0) * float64(res.Price)

		finalPrice = finalPrice + daysToAddVariation*percentageIncrease

		sumDaysToAddVariation = sumDaysToAddVariation + daysToAddVariation

		daysToAddVariations = append(daysToAddVariations, int(daysToAddVariation))
	}

	finalPrice = finalPrice + math.Abs(resDays-sumDaysToAddVariation)*float64(res.Price)

	return finalPrice, daysToAddVariations

}

func GetNumberOfDays(startDate time.Time, endDate time.Time) float64 {

	return (endDate.Sub(startDate).Hours() + 24) / 24
}

func (r *ReservationsHandler) CreateAvailability(c *gin.Context) {

	var av domain.Availability

	if err := c.ShouldBindJSON(&av); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := r.repo.InsertAvailability(&av)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		e := json.NewEncoder(c.Writer)
		e.Encode(message)
	}
	c.JSON(http.StatusOK, message)
}

func (r *ReservationsHandler) DeleteAvailability(c *gin.Context) {

	var av domain.Availability

	if err := c.ShouldBindJSON(&av); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mess := r.repo.DeleteAvailability(&av)
	if mess != "" {
		r.logger.Print("Database exception: ", mess)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *ReservationsHandler) DeletePriceVariation(c *gin.Context) {

	var pv domain.PriceVariation

	if err := c.ShouldBindJSON(&pv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mess := r.repo.DeletePriceVariation(&pv)
	if mess != "" {
		r.logger.Print("Database exception: ", mess)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *ReservationsHandler) GetAvailability(c *gin.Context) {

	location := c.Param("location")
	accommId := c.Param("accomm_id")

	availability, err := r.repo.GetAvailability(location, accommId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		e := json.NewEncoder(c.Writer)
		e.Encode(err.Error())
	}

	if availability == nil {
		return
	}

	e := json.NewEncoder(c.Writer)
	e.Encode(availability)
}

func (r *ReservationsHandler) CORSMiddleware() gin.HandlerFunc {
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

func (a *ReservationsHandler) Authorize(role string) gin.HandlerFunc { // Check if user is authorized as Host or Guest depending on the 'role' parameter
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
