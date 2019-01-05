package models

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	Role              string        `bson:"role" json:"role"`
	FirstName         string        `bson:"firstName" json:"firstName"`
	LastName          string        `bson:"lastName" json:"lastName"`
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
	Site              *Site         `bson:"-" json:"site,omitempty"`
	Create            time.Time     `bson:"create" json:"create"`
	LastUpdate        time.Time     `bson:"lastUpdate" json:"lastUpdate"`
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) GetFullname() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsFullnameEmpty() bool {
	return u.FirstName == "" || u.LastName == ""
}

func (u *User) ValidateRole() error {
	if u.Role != "corporate" && u.Role != "individual" && u.Role != "admin" {
		return utils.ErrInvalidUserRole
	}
	return nil
}

func (u *User) ValidateVat() error {
	if u.Vat != "N" && u.Vat != "Y" {
		return utils.ErrInvalidUserVat
	}
	return nil
}
