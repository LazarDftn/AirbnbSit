package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type AvailabilityByUser struct {
	AvailabilityID gocql.UUID `json:"availabilityId"`
	AccommID       string     `json:"accommId" bson:"accommId"`
	Name           string     `json:"name" bson:"name"`
	Location       string     `json:"location" bson:"location"`
	MinCapacity    int        `json:"minCapacity" bson:"minCapacity"`
	MaxCapacity    int        `json:"maxCapacity" bson:"maxCapacity"`
	StartDate      time.Time  `json:"startDate" bson:"startDate"`
	EndDate        time.Time  `json:"endDate" bson:"endDate"`
	UserId         string     `json:"userId" bson:"userId"`
}
