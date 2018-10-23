package income

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
)

func TestAddIncome(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("AddIncome", mock.AnythingOfType("*models.Income")).Return(&mocks.MockIncome, nil)

	uc := newUsecase(mockRepo)
	u, err := uc.AddIncome(&mocks.MockIncome)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockIncome.ID, u.ID)
	mockRepo.AssertExpectations(t)
}
