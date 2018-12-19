package invoice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{repo}
}

func (u *usecase) Create(i *models.Invoice) (*models.Invoice, error) {
	return u.repo.Create(i)
}
