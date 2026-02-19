package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type ForControllingUserIncome interface {
	GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error)
	AddIncome(u *models.Income) error
}

type ForGettingUserByID interface {
	GetByID(id string) (*models.User, error)
}
