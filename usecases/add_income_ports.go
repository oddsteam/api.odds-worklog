package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/models"

type ForUsingAddIncome interface {
	AddIncome(req *models.IncomeReq, uid string) (*models.Income, error)
}
