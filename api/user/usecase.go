package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo Repository
}

func newUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) CreateUser(m *models.User) (*models.User, error) {
	user, err := u.repo.CreateUser(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) GetUser() ([]*models.User, error) {
	users, err := u.repo.GetUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *usecase) GetUserByID(id string) (*models.User, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) UpdateUser(m *models.User) (*models.User, error) {
	user, err := u.repo.UpdateUser(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}
