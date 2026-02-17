package models

import (
	"github.com/globalsign/mgo/bson"
)

type Site struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Color string        `bson:"color" json:"color"`
}
