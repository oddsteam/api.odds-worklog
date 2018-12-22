package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Invoice struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	PoID       string        `bson:"poId" json:"poId"`
	InvoiceNo  string        `bson:"invoiceNo" json:"invoiceNo"`
	Amount     string        `bson:"amount" json:"amount"`
	Create     time.Time     `bson:"create" json:"create"`
	LastUpdate time.Time     `bson:"lastUpdate" json:"lastUpdate"`
}

type InvoiceNoRes struct {
	InvoiceNo string `json:"invoiceNo"`
}

type InvoiceNoReq struct {
	PoID string `json:"poId"`
}
