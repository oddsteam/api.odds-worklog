package invoice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase interface {
	Create(i *models.Invoice) (*models.Invoice, error)
	Get() ([]*models.Invoice, error)
	NextNo(id string) (string, error)
}

type Repository interface {
	Create(i *models.Invoice) (*models.Invoice, error)
	Get() ([]*models.Invoice, error)
	Last(id string) (*models.Invoice, error)
}
