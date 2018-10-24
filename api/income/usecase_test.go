package income

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
	userMocks "gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
)

func TestCalVAT(t *testing.T) {
	vat, err := calVAT("100000")
	assert.NoError(t, err)
	assert.Equal(t, "7000.000000", vat)

	vat, err = calVAT("123456")
	assert.NoError(t, err)
	assert.Equal(t, "8641.920000", vat)

	vat, err = calVAT("99999")
	assert.NoError(t, err)
	assert.Equal(t, "6999.930000", vat)
}

func TestCalWHT(t *testing.T) {
	vat, err := calWHT("100000")
	assert.NoError(t, err)
	assert.Equal(t, "3000.000000", vat)

	vat, err = calWHT("123456")
	assert.NoError(t, err)
	assert.Equal(t, "3703.680000", vat)

	vat, err = calWHT("99999")
	assert.NoError(t, err)
	assert.Equal(t, "2999.970000", vat)
}

func TestUsecaseAddIncome(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("AddIncome", mock.AnythingOfType("*models.Income")).Return(nil)

	mockUserRepo := new(userMocks.Repository)
	mockUserRepo.On("GetUserByID", mocks.MockIncome.UserID).Return(&userMocks.MockUserById, nil)

	uc := newUsecase(mockRepo, mockUserRepo)
	err := uc.AddIncome(&mocks.MockIncomeReq, mocks.MockIncome.UserID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
