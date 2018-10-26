package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase struct {
	mock.Mock
}

func (m *Usecase) AddIncome(income *models.IncomeReq, user *models.User) (*models.Income, error) {
	ret := m.Called(income, user)

	var r0 *models.Income
	if rf, ok := ret.Get(0).(func(*models.IncomeReq, *models.User) *models.Income); ok {
		r0 = rf(income, user)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Income)
	}

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.IncomeReq, *models.User) error); ok {
		r1 = rf(income, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error) {
	ret := m.Called(id, req, user)

	var r0 *models.Income
	if rf, ok := ret.Get(0).(func(string, *models.IncomeReq, *models.User) *models.Income); ok {
		r0 = rf(id, req, user)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Income)
	}

	var r1 error
	if rf, ok := ret.Get(0).(func(string, *models.IncomeReq, *models.User) error); ok {
		r1 = rf(id, req, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) GetIncomeStatusList() ([]*models.IncomeRes, error) {
	ret := m.Called()

	var r0 []*models.IncomeRes
	if rf, ok := ret.Get(0).(func() []*models.IncomeRes); ok {
		r0 = rf()
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*models.IncomeRes)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
