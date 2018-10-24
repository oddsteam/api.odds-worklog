package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase struct {
	mock.Mock
}

func (m *Usecase) AddIncome(income *models.IncomeReq, id string) (*models.IncomeRes, error) {
	ret := m.Called(income, id)

	var r0 *models.IncomeRes
	if rf, ok := ret.Get(0).(func(*models.IncomeReq, string) *models.IncomeRes); ok {
		r0 = rf(income, id)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.IncomeRes)
	}

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.IncomeReq, string) error); ok {
		r1 = rf(income, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
