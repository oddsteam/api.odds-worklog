package backoffice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	Get() ([]*models.UserIncome, error)
}

type Usecase interface {
	Get() ([]*models.UserIncome, error)
}
