package user

import "gitlab.odds.team/worklog/api.odds-worklog/models"

type Repository interface {
	createUser(u *models.User) (*models.User, error)
	getUser() ([]*models.User, error)
	getUserByID(id string) (*models.User, error)
	updateUser(u *models.User) (*models.User, error)
	// delete(id int64) (bool, error)
}

type Usecase interface {
	createUser(u *models.User) (*models.User, error)
	getUser() ([]*models.User, error)
	getUserByID(id string) (*models.User, error)
	updateUser(u *models.User) (*models.User, error)
	// delete(id int64) (bool, error)
}
