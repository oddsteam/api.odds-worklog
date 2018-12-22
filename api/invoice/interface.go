package invoice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase interface {
	Create(i *models.Invoice) (*models.Invoice, error)
	Get() ([]*models.Invoice, error)
	GetByPO(id string) ([]*models.Invoice, error)
	GetByID(id string) (*models.Invoice, error)
	NextNo(id string) (string, error)
	Update(i *models.Invoice) (*models.Invoice, error)
	Delete(id string) error
}

type Repository interface {
	Create(i *models.Invoice) (*models.Invoice, error)
	Get() ([]*models.Invoice, error)
	GetByPO(id string) ([]*models.Invoice, error)
	GetByID(id string) (*models.Invoice, error)
	Last(id string) (*models.Invoice, error)
	Update(i *models.Invoice) (*models.Invoice, error)
	Delete(id string) error
}
