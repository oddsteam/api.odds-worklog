package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Income struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	UserID      string        `bson:"userId" json:"userId"`
	TotalIncome string        `bson:"totalIncome" json:"totalIncome"`
	NetIncome   string        `bson:"netIncome" json:"netIncome"`
	SubmitDate  time.Time     `bson:"submitDate" json:"submitDate"`
	Note        string        `bson:"note" json:"note"`
	VAT         string        `bson:"vat" json:"vat"`
	WHT         string        `bson:"wht" json:"wht"`
}

type IncomeStatus struct {
	User       *User  `bson:"user" json:"user,omitempty"`
	SubmitDate string `bson:"submitDate" json:"submitDate"`
	Status     string `bson:"status" json:"status"`
}

type IncomeReq struct {
	TotalIncome string `bson:"totalIncome" json:"totalIncome"`
	Note        string `bson:"note" json:"note"`
}
