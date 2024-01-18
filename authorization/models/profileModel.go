package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID            primitive.ObjectID `bson:"_id",omitempty`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=10"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=20"`
	Username      *string            `json:"username"`
	Email         *string            `json:"email" validate:"required"`
	Address       *string            `json:"address" validate:"required"`
	Token         *string            `json:"token"`
	User_type     *string            `json:"user_type" validate:"required,eq=HOST|eq=GUEST"`
	Refresh_token *string            `json:"refresh_token"`
	Is_verified   bool               `json:"is_verified"`
}
