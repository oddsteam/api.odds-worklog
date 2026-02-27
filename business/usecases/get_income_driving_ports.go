package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForUsingGetIncome interface {
	GetIncomeByCurrentMonth(userId string) (*models.Income, error)
	GetIncomeByAllMonth(userId string) ([]*models.Income, error)
}
