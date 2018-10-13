package user

import "gitlab.odds.team/worklog/api.odds-worklog/models"

type Repository interface {
	CreateUser(u *models.User) (*models.User, error)
	GetUser() ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(u *models.User) (*models.User, error)
	DeleteUser(id string) error
	Login(authen *models.Login) (*models.Token, error)
}

type Usecase interface {
	CreateUser(u *models.User) (*models.User, error)
	GetUser() ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(u *models.User) (*models.User, error)
	DeleteUser(id string) error
	Login(authen *models.Login) (*models.Token, error)
}
