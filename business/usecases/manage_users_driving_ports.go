package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForUsingUsers interface {
	Create(u *models.User) (*models.User, error)
	Get() ([]*models.User, error)
	GetByRole(role string) ([]*models.User, error)
	GetByID(id string) (*models.User, error)
	GetBySiteID(id string) ([]*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(u *models.User, actor *models.UserClaims) (*models.User, error)
	Delete(id string) error
	UpdateStatusTavi(m []*models.StatusTavi, isAdmin bool) ([]*models.User, error)
}
