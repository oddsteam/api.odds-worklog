package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Export struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Filename string        `bson:"filename" json:"filename"`
	Date     time.Time     `bson:"date" json:"date"`
}
