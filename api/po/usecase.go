package po

import (
	"gitlab.odds.team/worklog/api.odds-worklog/api/customer"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo     Repository
	custRepo customer.Repository
}

func NewUsecase(r Repository, custRepo customer.Repository) Usecase {
	return &usecase{r, custRepo}
}

func (u *usecase) Create(m *models.Po) (*models.Po, error) {
	return u.repo.Create(m)
}

func (u *usecase) Update(m *models.Po) (*models.Po, error) {
	return u.repo.Update(m)
}

func (u *usecase) Get() ([]*models.Po, error) {
	return u.repo.Get()
}

func (u *usecase) GetByID(id string) (*models.Po, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) GetByCusID(id string) ([]*models.Po, error) {
	return u.repo.GetByCusID(id)
}
