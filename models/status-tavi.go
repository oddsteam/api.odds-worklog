package models

import (
	"github.com/globalsign/mgo/bson"
)

type StatusTavi struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	User *User         `bson:"user" json:"user,omitempty"`
}
