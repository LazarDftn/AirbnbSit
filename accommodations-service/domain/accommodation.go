package domain

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type Accommodation struct {
	Id          uuid.UUID
	Owner       string //change type string to User
	Name        string
	Location    string
	MinCapacity int
	MaxCapacity int
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
