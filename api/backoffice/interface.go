package backoffice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	Get() ([]*models.UserIncome, error)
	GetKey() (*models.BackOfficeKey, error)
}

type Usecase interface {
	Get() ([]*models.UserIncome, error)
	GetKey() (*models.BackOfficeKey, error)
}
