package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForUsingUpdateIncome interface {
	UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error)
}
