package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Income struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	UserID        string        `bson:"userId" json:"userId"`
	TotalIncome   string        `bson:"totalIncome" json:"totalIncome"`
	NetIncome     string        `bson:"netIncome" json:"netIncome"`
	SubmitDate    time.Time     `bson:"submitDate" json:"submitDate"`
	LastUpdate    time.Time     `bson:"lastUpdate" json:"lastUpdate"`
	Note          string        `bson:"note" json:"note"`
	VAT           string        `bson:"vat" json:"vat"`
	WHT           string        `bson:"wht" json:"wht"`
	WorkDate      string        `bson:"workDate" json:"workDate"`
	SpecialIncome string        `bson:"specialIncome" json:"specialIncome"`
}

type IncomeStatus struct {
	User       *User  `bson:"user" json:"user,omitempty"`
	SubmitDate string `bson:"submitDate" json:"submitDate"`
	Status     string `bson:"status" json:"status"`
	WorkDate   string `bson:"workDate" json:"workDate"`
}

type IncomeReq struct {
	SpecialIncome string `bson:"specialIncome" json:"specialIncome"`
	Note          string `bson:"note" json:"note"`
	WorkDate      string `bson:"workDate" json:"workDate"`
}
