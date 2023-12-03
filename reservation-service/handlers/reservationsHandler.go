package handlers

import (
	"fmt"
	"log"
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
