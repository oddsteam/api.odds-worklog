package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type UserIncome struct {
	ID                bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Role              string        `bson:"role" json:"role"`
	FirstName         string        `bson:"firstName" json:"firstName"`
	LastName          string        `bson:"lastName" json:"lastName"`
	CorporateName     string        `bson:"corporateName" json:"corporateName,omitempty"`
	Email             string        `bson:"email" json:"email"`
	BankAccountName   string        `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber string        `bson:"bankAccountNumber" json:"bankAccountNumber"`
	ThaiCitizenID     string        `bson:"thaiCitizenId" json:"thaiCitizenId,omitempty"`
	Vat               string        `bson:"vat" json:"vat,omitempty"`
	SlackAccount      string        `bson:"slackAccount" json:"slackAccount"`
	Transcript        string        `bson:"transcript" json:"transcript,omitempty"`
	SiteID            string        `bson:"siteId" json:"siteId,omitempty"`
	Project           string        `bson:"project" json:"project,omitempty"`
	ImageProfile      string        `bson:"imageProfile" json:"imageProfile,omitempty"`
	DegreeCertificate string        `bson:"degreeCertificate" json:"degreeCertificate,omitempty"`
	IDCard            string        `bson:"idCard" json:"idCard,omitempty"`
	Site              *Site         `bson:"-" json:"site,omitempty"`
	Create            time.Time     `bson:"create" json:"create"`
	LastUpdate        time.Time     `bson:"lastUpdate" json:"lastUpdate"`
	DailyIncome       string        `bson:"dailyIncome" json:"dailyIncome,omitempty"`
	Address           string        `bson:"address" json:"address,omitempty"`
	StatusTavi        bool          `bson:"statusTavi" json:"statusTavi"`
	Incomes           []*Income		`bson:"incomes" json:"incomes"`
}