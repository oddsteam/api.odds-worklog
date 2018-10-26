package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	AddIncome(u *models.Income) error
	GetIncomeUserNow(id, month string) (*models.Income, error)
	GetIncomeByID(incID, uID string) (*models.Income, error)
	UpdateIncome(income *models.Income) error
}

type Usecase interface {
	AddIncome(req *models.IncomeReq, user *models.User) (*models.Income, error)
	UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error)
	GetIncomeStatusList() ([]*models.IncomeRes, error)
}
