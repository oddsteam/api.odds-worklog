package models

import "github.com/globalsign/mgo/bson"

type StudentLoanList struct {
	List []StudentLoan `bson:"list"`
}

type StudentLoan struct {
	ID       bson.ObjectId `bson:"_id"`
	Fullname string        `bson:"customerName"`
	Amount   int           `bson:"paidAmount"`
}
