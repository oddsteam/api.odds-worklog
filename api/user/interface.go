package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	CreateUser(u *models.User) (*models.User, error)
	GetUser() ([]*models.User, error)
	GetUserByRole(role string) ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserBySiteID(id string) ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(u *models.User) (*models.User, error)
	DeleteUser(id string) error
}

type Usecase interface {
	CreateUser(u *models.User) (*models.User, error)
	GetUser() ([]*models.User, error)
	GetUserByRole(role string) ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserBySiteID(id string) ([]*models.User, error)
	UpdateUser(u *models.User, isAdmin bool) (*models.User, error)
	DeleteUser(id string) error
}
