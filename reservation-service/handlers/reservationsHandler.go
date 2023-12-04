package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"reservation-service/domain"
	"reservation-service/repositories"

	"github.com/gin-gonic/gin"
)

type KeyProduct struct{}

type ReservationsHandler struct {
	logger *log.Logger
	// NoSQL: injecting reservation repository
	repo *repositories.ReservationRepo
}

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
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

func (rr *ReservationsHandler) GetResByAccommodation(c *gin.Context) {

	accommId := c.Param("accomm_id")

	reservations, err := rr.repo.GetReservationsByAccomm(accommId)
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
		/*price,*/ price, days := CompareAndCalculate(variations[0], res)
		fmt.Print(price)
		fmt.Print(days)
		retVal := struct {
			Price       int
			Percentage  int
			OverlapDays int
		}{
			Price:       int(price),
			Percentage:  variations[0].Percentage,
			OverlapDays: int(days),
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

	err := rr.repo.InsertPriceVariation(&vr)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

func (rr *ReservationsHandler) GetPriceVariationByAccommId(c *gin.Context) {
	accommId := c.Param("accomm_id")

	variations, err := rr.repo.GetVariationsByAccommId(accommId)
	if err != nil {
		rr.logger.Print("Database exception: ", err)
	}

	if variations == nil {
		return
	}

	c.JSON(http.StatusOK, variations)
}

func CompareAndCalculate(vr domain.PriceVariation, res domain.Reservation) (float64, float64) {

	resDays := (res.EndDate.Sub(res.StartDate).Hours() + 24) / 24

	daysToAddVariation := math.Min(float64(res.EndDate.Day()), float64(vr.EndDate.Day())) -
		math.Max(float64(res.StartDate.Day()), float64(vr.StartDate.Day())) + 1

	percentageIncrease := (float64(100+vr.Percentage) / 100.0) * float64(res.Price)

	finalPrice := daysToAddVariation*percentageIncrease + math.Abs(resDays-daysToAddVariation)*float64(res.Price)

	return finalPrice, daysToAddVariation

}

func (a *ReservationsHandler) CORSMiddleware() gin.HandlerFunc {
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
