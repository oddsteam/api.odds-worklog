package models

import (
	"github.com/globalsign/mgo/bson"
)

type BackOfficeKey struct {
	ID                bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Key              string        `bson:"key" json:"key"`
}