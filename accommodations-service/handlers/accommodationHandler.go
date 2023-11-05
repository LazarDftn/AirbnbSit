package handlers

import (
	"accommodations-service/repositories"
	"log"
	"net/http"
)

type AccommodationsHandler struct {
	logger *log.Logger
	repo   *repositories.AccommodationRepo
}

func NewAccommodationsHandler(l *log.Logger, r *repositories.AccommodationRepo) *AccommodationsHandler {
	return &AccommodationsHandler{l, r}
}

func (a *AccommodationsHandler) GetAllPatients(rw http.ResponseWriter, h *http.Request) {
	patients, err := a.repo.GetAll()
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if patients == nil {
		return
	}

	err = patients.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}
