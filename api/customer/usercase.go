package customer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) Create(m *models.Customer) (*models.Customer, error) {
	return u.repo.Create(m)
}

func (u *usecase) Get() ([]*models.Customer, error) {
	return u.repo.Get()
}

func (u *usecase) GetByID(id string) (*models.Customer, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) Update(m *models.Customer) (*models.Customer, error) {
	customer, err := u.repo.GetByID(m.ID.Hex())
	if err != nil {
		return nil, err
	}

	if m.Name != "" {
		customer.Name = m.Name
	}

	if m.Address != "" {
		customer.Address = m.Address
	}

	customer, err = u.repo.Update(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (u *usecase) Delete(id string) error {
	return u.repo.Delete(id)
}
