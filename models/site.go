package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Site struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Color string        `bson:"color" json:"color"`
}
