package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type Availability struct {
	AvailabilityID gocql.UUID `json:"availabilityId"`
	AccommID       string     `json:"accommId"`
	StartDate      time.Time  `json:"startDate"`
	EndDate        time.Time  `json:"endDate"`
}
