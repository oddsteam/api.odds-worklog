package po

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) Create(m *models.Po) (*models.Po, error) {
	if m.CustomerId == "" {
		return nil, utils.ErrEmptyCustomerId
	}
	return u.repo.Create(m)
}

func (u *usecase) Update(m *models.Po) (*models.Po, error) {
	return u.repo.Update(m)
}
