package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) AddIncome(_u *models.Income) error {
	ret := m.Called(_u)

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Income) error); ok {
		r1 = rf(_u)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}

func (m *Repository) GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error) {
	ret := m.Called(id, fromYear, fromMonth)

	var r0 *models.Income
	if rf, ok := ret.Get(0).(func(string, int, time.Month) *models.Income); ok {
		r0 = rf(id, fromYear, fromMonth)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Income)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, time.Month) error); ok {
		r1 = rf(id, fromYear, fromMonth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Repository) GetIncomeByID(id1, id2 string) (*models.Income, error) {
	ret := m.Called(id1, id2)

	var r0 *models.Income
	if rf, ok := ret.Get(0).(func(string, string) *models.Income); ok {
		r0 = rf(id1, id2)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Income)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id1, id2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *Repository) UpdateIncome(income *models.Income) error {
	ret := m.Called(income)

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Income) error); ok {
		r1 = rf(income)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}

func (m *Repository) AddExport(ep *models.Export) error {
	ret := m.Called(ep)

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Export) error); ok {
		r1 = rf(ep)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}
