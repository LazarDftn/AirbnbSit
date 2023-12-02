package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type Reservation struct {
	ReservationID gocql.UUID `json:"reservationId"`
	AccommID      string     `json:"accommId"`
	GuestEmail    string     `json:"guestEmail"`
	Price         int        `json:"price"`
	NumOfPeople   int        `json:"numOfPeople"`
	StartDate     time.Time  `json:"startDate"`
	EndDate       time.Time  `json:"endDate"`
}
