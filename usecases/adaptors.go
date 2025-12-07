package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type ForGettingIncomeData interface {
	GetAllIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
	GetStudentLoans() models.StudentLoanList
}

type ForControllingIncomeData interface {
	AddExport(ep *models.Export) error
}
