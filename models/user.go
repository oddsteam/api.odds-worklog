package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	FullNameTh        string        `bson:"fullnameTh" json:"fullnameTh"`
	FullNameEn        string        `bson:"fullnameEn" json:"fullnameEn"`
	Email             string        `bson:"email" json:"email"`
	BankAccountName   string        `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber string        `bson:"bankAccountNumber" json:"bankAccountNumber"`
	ThaiCitizenID     string        `bson:"thaiCitizenId" json:"thaiCitizenId,omitempty"`
	CorporateFlag     string        `bson:"corporateFlag" json:"corporateFlag"`
}
