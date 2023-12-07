package domain

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gocql/gocql"
)

type Reservation struct {
	ReservationID gocql.UUID `json:"reservationId"`
	Location      string     `json:"location"`
	AccommID      string     `json:"accommId"`
	GuestEmail    string     `json:"guestEmail"`
	HostEmail     string     `json:"hostEmail"`
	Price         int        `json:"price"`
	NumOfPeople   int        `json:"numOfPeople"`
	StartDate     time.Time  `json:"startDate"`
	EndDate       time.Time  `json:"endDate"`
}

type Reservations []*Reservation

func (a *Reservations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Reservation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Reservation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
