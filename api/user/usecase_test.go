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
	assert.Equal(t, mocks.MockUsers[0].FullNameTh, u[0].FullNameTh)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetUserByID", "1234567890").Return(&mocks.MockUserById, nil)

	uc := newUsecase(mockRepo)
	u, err := uc.GetUserByID(string(mocks.MockUserById.ID))
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockUserById.FullNameTh, u.FullNameTh)
	mockRepo.AssertExpectations(t)
}
func TestDeleteUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("DeleteUser", "1234567890").Return(nil)

	uc := newUsecase(mockRepo)
	u := uc.DeleteUser(string(mocks.MockUserById.ID))

	assert.Equal(t, nil, u)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUserById, nil)
	uc := newUsecase(mockRepo)
	u, err := uc.UpdateUser(&mocks.MockUserById)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.MockUser.FullNameTh, u.FullNameTh)
	mockRepo.AssertExpectations(t)
}
