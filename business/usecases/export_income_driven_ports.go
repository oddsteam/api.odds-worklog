package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type ForGettingIncomeDataInTheMonth interface {
	GetAllIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
}

type ForListStudentLoansInTheMonth interface {
	GetStudentLoans() models.StudentLoanList
}

type ForLoggingExport interface {
	AddExport(ep *models.Export) error
}

type ForLoggingSAPExportFailure interface {
	LogSAPExportFailure(log *models.SAPExportFailureLog) error
}

type ForWritingCSVFile interface {
	WriteFile(name string, ics models.PayrollCycle) (string, error)
}

type ForWritingSAPFile interface {
	WriteFile(name string, ics models.PayrollCycle, dateEff time.Time) (string, error)
}
