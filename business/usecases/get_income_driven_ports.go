package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type ForReadingUserIncome interface {
	GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error)
	GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error)
}
