package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Export struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	Filename      string        `bson:"filename" json:"filename"`
	CorporateFlag string        `bson:"corporateFlag" json:"corporateFlag"`
	Date          time.Time     `bson:"date" json:"date"`
}
