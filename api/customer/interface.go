package customer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	CreateCustomer(u *models.Customer) (*models.Customer, error)
	GetCustomers() ([]*models.Customer, error)
	GetCustomerByID(id string) (*models.Customer, error)
	UpdateCustomer(custmoer *models.Customer) (*models.Customer, error)
	DeleteCustomer(id string) error
}

type Usecase interface {
	CreateCustomer(u *models.Customer) (*models.Customer, error)
	GetCustomers() ([]*models.Customer, error)
	GetCustomerByID(id string) (*models.Customer, error)
	UpdateCustomer(m *models.Customer, isAdmin bool) (*models.Customer, error)
	DeleteCustomer(id string) error
}
