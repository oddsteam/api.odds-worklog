package models

import "gopkg.in/mgo.v2/bson"

type Income struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	UserID      string        `bson:"userId" json:"userId"`
	TotalIncome string        `bson:"totalIncome" json:"totalIncome"`
	SubmitDate  string        `bson:"submitDate" json:"submitDate"`
	Note        string        `bson:"note" json:"note"`
	VAT         string        `bson:"vat" json:"vat"`
	WHT         string        `bson:"whr" json:"wht"`
}
