package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUser, nil)

	uc := newUsecase(mockRepo)
	u, err := uc.CreateUser(&mocks.MockUser)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockUser.ID, u.ID)
	mockRepo.AssertExpectations(t)
}
