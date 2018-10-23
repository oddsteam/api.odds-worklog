package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	FullName          string        `bson:"fullname" json:"fullname"`
	Email             string        `bson:"email" json:"email"`
	BankAccountName   string        `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber string        `bson:"bankAccountNumber" json:"bankAccountNumber"`
	ThaiCitizenID     string        `bson:"thaiCitizenId" json:"thaiCitizenId"`
	CoperateFlag      string        `bson:"coperateFlag" json:"coperateFlag"`
}
