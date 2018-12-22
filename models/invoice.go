package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Invoice struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	PoID      string        `bson:"poId" json:"poId"`
	InvoiceNo string        `bson:"invoiceNo" json:"invoiceNo"`
	Amount    string        `bson:"amount" json:"amount"`
}
