package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForUsingAddIncome interface {
	AddIncome(req *models.IncomeReq, uid string) (*models.Income, error)
}
