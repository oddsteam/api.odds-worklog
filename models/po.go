package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Po struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	CustomerId string        `bson:"customerId" json:"customerId"`
	Name       string        `bson:"name" json:"name"`
	Amount     string        `bson:"amount" json:"amount"`
	Create     time.Time     `bson:"create" json:"create"`
	LastUpdate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
}
