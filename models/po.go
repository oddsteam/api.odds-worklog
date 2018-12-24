package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Po struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	CustomerId string        `bson:"customerId" json:"customerId"`
	Name       string        `bson:"name" json:"name"`
	Create     time.Time     `bson:"create" json:"create"`
	LastUpdate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
}
