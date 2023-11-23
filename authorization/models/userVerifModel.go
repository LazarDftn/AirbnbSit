package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserVerifModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VerifUsername *string            `bson:"verifUsername" json:"verifUsername"`
	Code          *string            `bson:"code" json:"code"`
	Created_at    *time.Time         `bson:"created_at,omitempty" json:"created_at"`
}
