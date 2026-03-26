package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// SAPExportFailureLog is stored when SAP file export fails (e.g. Windows-874 encoding error) for operations debugging.
type SAPExportFailureLog struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	CreatedAt       time.Time     `bson:"createdAt" json:"createdAt"`
	Role            string        `bson:"role" json:"role"`
	StartDate       time.Time     `bson:"startDate" json:"startDate"`
	EndDate         time.Time     `bson:"endDate" json:"endDate"`
	DateEffective   time.Time     `bson:"dateEffective" json:"dateEffective"`
	IncomeID        string        `bson:"incomeId" json:"incomeId"`
	UserID          string        `bson:"userId" json:"userId"`
	BankAccountName string        `bson:"bankAccountName" json:"bankAccountName"`
	RowIndex        int           `bson:"rowIndex" json:"rowIndex"`
	LineKind        string        `bson:"lineKind" json:"lineKind"`
	ErrorMessage    string        `bson:"errorMessage" json:"errorMessage"`
	UnderlyingError string        `bson:"underlyingError,omitempty" json:"underlyingError,omitempty"`
}
