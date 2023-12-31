package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type Availability struct {
	AvailabilityID gocql.UUID `json:"availabilityId"`
	AccommID       string     `json:"accommId"`
	Name           string     `json:"name"`
	Location       string     `json:"location"`
	MinCapacity    int        `json:"minCapacity"`
	MaxCapacity    int        `json:"maxCapacity"`
	StartDate      time.Time  `json:"startDate"`
	EndDate        time.Time  `json:"endDate"`
}
