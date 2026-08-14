package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForManagingUsers interface {
	Create(u *models.User) (*models.User, error)
	CreateArchivedUser(a models.User) (*models.ArchivedUser, error)
	Get() ([]*models.User, error)
	GetByRole(role string) ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetBySiteID(id string) ([]*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User) (*models.User, error)
	Delete(id string) error
}
