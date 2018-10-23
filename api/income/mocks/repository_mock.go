package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) AddIncome(_u *models.Income) (*models.Income, error) {
	ret := m.Called(_u)

	var r0 *models.Income
	if rf, ok := ret.Get(0).(func(*models.Income) *models.Income); ok {
		r0 = rf(_u)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.Income)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Income) error); ok {
		r1 = rf(_u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
