package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner       string             `bson:"owner" json:"owner"`
	OwnerId     string             `bson:"ownerId" json:"ownerId"`
	Name        string             `bson:"name" json:"name"`
	Location    string             `bson:"location" json:"location"`
	Benefits    string             `bson:"benefits" json:"benefits"`
	MinCapacity int                `bson:"minCapacity" json:"minCapacity"`
	MaxCapacity int                `bson:"maxCapacity" json:"maxCapacity"`
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
