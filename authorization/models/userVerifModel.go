package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserVerifModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VerifUsername *string            `bson:"verifUsername" json:"verifUsername"`
	Code          *string            `bson:"code" json:"code"`
}
