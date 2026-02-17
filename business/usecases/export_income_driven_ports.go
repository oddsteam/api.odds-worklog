package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type ForGettingIncomeData interface {
	GetAllIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
	GetStudentLoans() models.StudentLoanList
}

type ForControllingIncomeData interface {
	AddExport(ep *models.Export) error
}

type ForWritingCSVFile interface {
	WriteFile(name string, ics models.PayrollCycle) (string, error)
}

type ForWritingSAPFile interface {
	WriteFile(name string, ics models.PayrollCycle, dateEff time.Time) (string, error)
}
