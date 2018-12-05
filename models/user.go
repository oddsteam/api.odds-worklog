package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	FirstName         string        `bson:"firstName" json:"firstName"`
	LastName          string        `bson:"lastName" json:"lastName"`
	Email             string        `bson:"email" json:"email"`
	BankAccountName   string        `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber string        `bson:"bankAccountNumber" json:"bankAccountNumber"`
	ThaiCitizenID     string        `bson:"thaiCitizenId" json:"thaiCitizenId,omitempty"`
	CorporateFlag     string        `bson:"corporateFlag" json:"corporateFlag"`
	Vat               string        `bson:"vat" json:"vat,omitempty"`
	SlackAccount      string        `bson:"slackAccount" json:"slackAccount"`
	Transcript        string        `bson:"transcript" json:"transcript,omitempty"`
	SiteID            string        `bson:"siteId" json:"SiteId,omitempty"`
}

func (u *User) GetFullname() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsFullnameEmpty() bool {
	return u.FirstName == "" || u.LastName == ""
}
