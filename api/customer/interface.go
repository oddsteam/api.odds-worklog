package customer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	Create(u *models.Customer) (*models.Customer, error)
	Get() ([]*models.Customer, error)
	GetByID(id string) (*models.Customer, error)
	Update(custmoer *models.Customer) (*models.Customer, error)
	Delete(id string) error
}

type Usecase interface {
	Create(u *models.Customer) (*models.Customer, error)
	Get() ([]*models.Customer, error)
	GetByID(id string) (*models.Customer, error)
	Update(m *models.Customer) (*models.Customer, error)
	Delete(id string) error
}
