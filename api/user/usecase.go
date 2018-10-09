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

	m.ID = user.ID
	return m, nil
}

func (u *usecase) getUser() ([]*models.User, error) {
	users, err := u.repo.getUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}
