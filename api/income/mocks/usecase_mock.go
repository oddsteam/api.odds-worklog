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

func (m *Usecase) GetIncomeStatusList(corporateFlag string) ([]*models.IncomeStatus, error) {
	ret := m.Called(corporateFlag)

	var r0 []*models.IncomeStatus
	if rf, ok := ret.Get(0).(func(string) []*models.IncomeStatus); ok {
		r0 = rf(corporateFlag)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*models.IncomeStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(corporateFlag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) GetIncomeByUserIdAndCurrentMonth(userId string) (*models.Income, error) {
	ret := m.Called(userId)

	var r0 *models.Income
	if rf, ok := ret.Get(0).(func(string) *models.Income); ok {
		r0 = rf(userId)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Income)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) ExportIncome(corporateFlag string) (string, error) {
	ret := m.Called(corporateFlag)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(corporateFlag)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(corporateFlag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Usecase) DropIncome() error {
	ret := m.Called()

	var r1 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(0)
	}

	return r1
}
