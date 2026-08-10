package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type ForGettingIncomeFromTimesheetInTheMonth interface {
	GetAllByRoleStartDateAndEndDate(role string, startDate, endDate time.Time) ([]*models.IncomeFromTimesheet, error)
}
