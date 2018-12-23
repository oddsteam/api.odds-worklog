package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Customer struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Address    string        `bson:"address" json:"address"`
	Create     time.Time     `bson:"create" json:"create"`
	LastUpdate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
}
