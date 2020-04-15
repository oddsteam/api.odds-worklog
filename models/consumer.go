package models

import "github.com/globalsign/mgo/bson"

type Consumer struct {
	ID       bson.ObjectId `bson:"_id"`
	ClientID string        `bson:"clientId"`
}
