package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Site struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Color string        `bson:"color" json:"color"`
}
