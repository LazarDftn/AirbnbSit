package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type PriceVariation struct {
	VariationID gocql.UUID `json:"variationId"`
	Location    string     `json:"location"`
	AccommID    string     `json:"accommId"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     time.Time  `json:"endDate"`
	Percentage  int        `json:"percentage"`
}
