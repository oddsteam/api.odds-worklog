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

func (m *Usecase) GetUser() ([]*models.User, error) {
	ret := m.Called()

	var r0 []*models.User
	if rf, ok := ret.Get(0).(func() []*models.User); ok {
		r0 = rf()
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) GetUserByID(id string) (*models.User, error) {
	ret := m.Called(id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(id)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) UpdateUser(_u *models.User) (*models.User, error) {
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

func (m *Usecase) DeleteUser(id string) error {
	ret := m.Called(id)

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r1
}

func (m *Usecase) Login(_u *models.Login) (*models.Token, error) {
	ret := m.Called(_u)

	var r0 *models.Token
	if rf, ok := ret.Get(0).(func(*models.Login) *models.Token); ok {
		r0 = rf(_u)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Token)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Login) error); ok {
		r1 = rf(_u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
