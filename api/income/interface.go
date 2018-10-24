package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	AddIncome(u *models.Income) error
}

type Usecase interface {
	AddIncome(u *models.IncomeReq, id string) (*models.IncomeRes, error)
}
