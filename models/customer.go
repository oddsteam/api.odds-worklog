package models

import "gopkg.in/mgo.v2/bson"

type Customer struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Name    string        `bson:"name" json:"name"`
	Address string        `bson:"address" json:"address"`
}
