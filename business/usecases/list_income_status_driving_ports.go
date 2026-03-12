package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForUsingListIncomeStatus interface {
	GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error)
}
