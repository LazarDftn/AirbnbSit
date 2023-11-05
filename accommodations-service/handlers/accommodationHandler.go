package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/repositories"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger
	repo   *repositories.AccommodationRepo
}

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
	a.repo.Insert(accomm)
	c.Writer.WriteHeader(http.StatusCreated)
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

func (p *AccommodationsHandler) MiddlewareAccommodationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		accommodation := &domain.Accommodation{}
		err := accommodation.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, accommodation)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (p *AccommodationsHandler) MiddlewareContentTypeSet() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
	}
}
