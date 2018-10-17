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

func TestGetUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetUser").Return(mocks.MockUsers, nil)

	uc := newUsecase(mockRepo)
	u, err := uc.GetUser()
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockUsers[0].FullName, u[0].FullName)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetUserByID", "1234567890").Return(&mocks.MockUserById, nil)

	uc := newUsecase(mockRepo)
	u, err := uc.GetUserByID(string(mocks.MockUserById.ID))
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockUserById.FullName, u.FullName)
	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("Login", mock.AnythingOfType("*models.Login")).Return(&mocks.MockToken, nil)
	uc := newUsecase(mockRepo)
	u, err := uc.Login(&mocks.Login)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockToken.Token, u.Token)
	mockRepo.AssertExpectations(t)
}
