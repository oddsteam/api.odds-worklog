package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Export struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Filename string        `bson:"filename" json:"filename"`
	Date     time.Time     `bson:"date" json:"date"`
}
