package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type ForReadingIncomeFromTimesheetByUser interface {
	GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeFromTimesheet, error)
	GetByUserIdAllMonth(userId string) ([]*models.IncomeFromTimesheet, error)
}
