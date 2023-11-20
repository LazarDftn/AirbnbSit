package domain

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner         string             `bson:"owner" json:"owner"` //change type string to User
	Name          string             `bson:"name" json:"name"`
	Location      string             `bson:"location" json:"location"`
	Benefits      string             `bson:"benefits" json:"benefits"`
	MinCapacity   int                `bson:"minCapacity" json:"minCapacity"`
	MaxCapacity   int                `bson:"maxCapacity" json:"maxCapacity"`
	Price         int                `bson:"price" json:"price"`                 //cena
	DiscPrice     int                `bson:"discPrice" json:"discPrice"`         //cena sa popustom
	DiscDateStart time.Time          `bson:"discDateStart" json:"discDateStart"` //pocetak vremenskog perioda popusta
	DiscDateEnd   time.Time          `bson:"discDateEnd" json:"discDateEnd"`     //kraj vremenskog perioda popusta
	DiscWeekend   bool               `bson:"discWeekend" json:"discWeekend"`     //ponavljajuci popust vikendom
	PayPer        int                `bson:"payPer" json:"payPer"`               //nacin placanja, 0=ukupno, 1=po osobi
}

type Accommodations []*Accommodation

func (a *Accommodations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Accommodation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Accommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
