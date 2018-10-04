package service

import (
	"github.com/globalsign/mgo/bson"
)

var Config = struct {
	APPName string `default:"worklog"`
	HOST    string
	DB      struct {
		Host string
	}
}{}

type User struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	FullName          string        `bson:"fullname" json:"fullname"`
	Email             string        `bson:"email" json:"email"`
	BankAccountName   string        `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber string        `bson:"bankAccountNumber" json:"bankAccountNumber"`
	TotalIncome       string        `bson:"totalIncome" json:"totalIncome"`
	SubmitDate        string        `bson:"submitDate" json:"submitDate"`
	CardNumber        string        `bson:"cardNumber" json:"cardNumber"`
}
