package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name *string            `json:"first_name" validate:"required,min=2,max=10"`
	Last_name  *string            `json:"last_name" validate:"required,min=2,max=20"`
	Username   *string            `json:"username"`
	Email      *string            `json:"email" validate:"required"`
	Address    *string            `json:"address" validate:"required"`
	User_type  *string            `json:"user_type" validate:"required,eq=HOST|eq=GUEST"`
}

type Users []User

func (a *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
