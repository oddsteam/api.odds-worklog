package mocks

import (
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
