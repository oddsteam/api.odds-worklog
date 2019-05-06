package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	Create(u *models.User) (*models.User, error)
	Get() ([]*models.User, error)
	GetByRole(role string) ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetBySiteID(id string) ([]*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Delete(id string) error
}

type Usecase interface {
	Create(u *models.User) (*models.User, error)
	Get() ([]*models.User, error)
	GetByRole(role string) ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetBySiteID(id string) ([]*models.User, error)
	Update(u *models.User, isAdmin bool) (*models.User, error)
	Delete(id string) error
	GetByEmail(email string) (*models.User, error)
	UpdateStatusTavi(m *models.User, isAdmin bool) (*models.User, error)
}
