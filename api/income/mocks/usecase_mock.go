package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase struct {
	mock.Mock
}

func (m *Usecase) AddIncome(income *models.IncomeReq, id string) error {
	ret := m.Called(income, id)

	var r1 error
	if rf, ok := ret.Get(0).(func(*models.IncomeReq, string) error); ok {
		r1 = rf(income, id)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}
