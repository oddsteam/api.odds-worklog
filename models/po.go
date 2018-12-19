package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Po struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	CustomerId string `bson:"customerId json:"customerId"`
	Name string `bson:"name" json:"name"`   			
}