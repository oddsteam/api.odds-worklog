package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForUpdatingUserIncome interface {
	GetIncomeByID(incID, uID string) (*models.Income, error)
	UpdateIncome(income *models.Income) error
}
