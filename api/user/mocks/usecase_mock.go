package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase struct {
	mock.Mock
}

func (m *Usecase) CreateUser(_u *models.User) (*models.User, error) {
	ret := m.Called(_u)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(*models.User) *models.User); ok {
		r0 = rf(_u)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(_u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (u *Usecase) GetUser() ([]*models.User, error) {
	return nil, nil
}

func (u *Usecase) GetUserByID(id string) (*models.User, error) {
	return nil, nil
}

func (u *Usecase) UpdateUser(m *models.User) (*models.User, error) {
	return nil, nil
}

func (u *Usecase) DeleteUser(id string) error {
	return nil
}

func (u *Usecase) Login(user *models.Login) (*models.Token, error) {
	return nil, nil
}
