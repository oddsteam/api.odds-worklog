package customer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

func NewUsecase(r Repository, user user.Repository) Usecase {
	return &usecase{r, user}
}

func (u *usecase) CreateCustomer(m *models.Customer) (*models.Customer, error) {
	return u.repo.CreateCustomer(m)
}

func (u *usecase) GetCustomers() ([]*models.Customer, error) {
	customers, err := u.repo.GetCustomers()
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (u *usecase) GetCustomerByID(id string) (*models.Customer, error) {
	return u.repo.GetCustomerByID(id)
}

func (u *usecase) UpdateCustomer(m *models.Customer, isAdmin bool) (*models.Customer, error) {
	if !isAdmin {
		return nil, utils.ErrInvalidUserRole
	}

	customer, err := u.repo.GetCustomerByID(m.ID.Hex())
	if err != nil {
		return nil, err
	}

	if m.Name != "" {
		customer.Name = m.Name
	}

	if m.Address != "" {
		customer.Address = m.Address
	}

	customer, err = u.repo.UpdateCustomer(m)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (u *usecase) DeleteCustomer(id string) error {
	return u.repo.DeleteCustomer(id)
}
