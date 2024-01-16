package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCredentialsModel struct {
	ID            primitive.ObjectID `bson:"_id"`
	Email         *string            `json:"email" validate:"required"`
	Password      *string            `json:"Password" validate:"required,min=12"`
	Is_verified   bool               `json:"is_verified"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
}
