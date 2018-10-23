package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo Repository
}

func newUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) AddIncome(m *models.Income) (*models.Income, error) {
	user, err := u.repo.AddIncome(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}
