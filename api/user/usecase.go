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

func (u *usecase) createUser(m *models.User) (*models.User, error) {
	user, err := u.repo.createUser(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) getUser() ([]*models.User, error) {
	users, err := u.repo.getUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *usecase) getUserByID(id string) (*models.User, error) {
	user, err := u.repo.getUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) updateUser(m *models.User) (*models.User, error) {
	user, err := u.repo.updateUser(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) deleteUser(id string) error {
	return u.repo.deleteUser(id)
}
