package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Income struct {
	ID                     primitive.ObjectID `bson:"_id" json:"id"`
	UserID                 string             `bson:"userId" json:"userId"`
	TotalIncome            string             `bson:"totalIncome" json:"totalIncome"`
	NetIncome              string             `bson:"netIncome" json:"netIncome"`
	DailyIncomeBeforeTax   string             `bson:"dailyIncomeBeforeTax" json:"dailyIncomeBeforeTax"`
	SpecialIncomeBeforeTax string             `bson:"specialIncomeBeforeTax" json:"specialIncomeBeforeTax"`
	NetDailyIncome         string             `bson:"netDailyIncome" json:"netDailyIncome"`
	WorkDate               string             `bson:"workDate" json:"workDate"`
	SubmitDate             time.Time          `bson:"submitDate" json:"submitDate"`
	LastUpdate             time.Time          `bson:"lastUpdate" json:"lastUpdate"`
	Note                   string             `bson:"note" json:"note"`
	VAT                    string             `bson:"vat" json:"vat"`
	WHT                    string             `bson:"wht" json:"wht"`
	WHTRate                float64            `bson:"whtRate" json:"whtRate,omitempty"`
	SpecialIncome          string             `bson:"specialIncome" json:"specialIncome"`
	NetSpecialIncome       string             `bson:"netSpecialIncome" json:"netSpecialIncome"`
	WorkingHours           string             `bson:"workingHours" json:"workingHours"`
	ExportStatus           bool               `bson:"exportStatus" json:"exportStatus"`
	ThaiCitizenID          string             `bson:"thaiCitizenId" json:"thaiCitizenId,omitempty"`
	Name                   string             `bson:"name" json:"name,omitempty"`
	BankAccountName        string             `bson:"bankAccountName" json:"bankAccountName"`
	BankAccountNumber      string             `bson:"bankAccountNumber" json:"bankAccountNumber"`
	Email                  string             `bson:"email" json:"email"`
	Phone                  string             `bson:"phone" json:"phone"`
	DailyRate              float64            `bson:"dailyRate"`
	IsVATRegistered        bool               `bson:"isVATRegistered"`
	Role                   string             `bson:"role" json:"role"`
	PeakCode               string             `bson:"peakCode" json:"peakCode,omitempty"`
	// SiteName is filled at export time (not persisted).
	SiteName string `bson:"-" json:"-"`
}

func (i *Income) EffectiveWHTRate() float64 {
	if i == nil || i.WHTRate == 0 {
		return LegacyWHTRate
	}
	return i.WHTRate
}

func (i *Income) WHTRateForExport() string {
	return fmt.Sprintf("%.0f%%", i.EffectiveWHTRate()*100)
}

type IncomeStatus struct {
	User         *User  `bson:"user" json:"user,omitempty"`
	SubmitDate   string `bson:"submitDate" json:"submitDate"`
	Status       string `bson:"status" json:"status"`
	WorkDate     string `bson:"workDate" json:"workDate"`
	WorkingHours string `bson:"workingHours" json:"workingHours"`
}
