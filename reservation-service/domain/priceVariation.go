package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type PriceVariation struct {
	VariationID gocql.UUID `json:"variationId"`
	AccommID    string     `json:"accommId"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     time.Time  `json:"endDate"`
	Percentage  int        `json:"percentage"`
}
