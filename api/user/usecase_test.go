package user

import (
	"testing"

	mock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &mock.MockUser

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().CreateUser(user).Return(user, nil)
	mockRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil)

	uc := NewUsecase(mockRepo)
	userRes, err := uc.CreateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, userRes)
	assert.Equal(t, user.ID, userRes.ID)
}

// func TestGetUser(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("GetUser").Return(mocks.MockUsers, nil)

// 	uc := NewUsecase(mockRepo)
// 	u, err := uc.GetUser()
// 	assert.NoError(t, err)
// 	assert.NotNil(t, u)
// 	assert.Equal(t, mocks.MockUsers[0].FullNameEn, u[0].FullNameEn)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUserByType(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("GetUserByType", "Y").Return(mocks.MockUsers, nil)

// 	uc := NewUsecase(mockRepo)
// 	list, err := uc.GetUserByType("Y")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, list)
// 	assert.Equal(t, mocks.MockUsers, list)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUserByID(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("GetUserByID", "1234567890").Return(&mocks.MockUserById, nil)

// 	uc := NewUsecase(mockRepo)
// 	u, err := uc.GetUserByID(string(mocks.MockUserById.ID))
// 	assert.NoError(t, err)
// 	assert.NotNil(t, u)
// 	assert.Equal(t, mocks.MockUserById.FullNameEn, u.FullNameEn)
// 	mockRepo.AssertExpectations(t)
// }
// func TestDeleteUser(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("DeleteUser", "1234567890").Return(nil)

// 	uc := NewUsecase(mockRepo)
// 	u := uc.DeleteUser(string(mocks.MockUserById.ID))

// 	assert.Equal(t, nil, u)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUpdateUser(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUserById, nil)
// 	uc := NewUsecase(mockRepo)
// 	u, err := uc.UpdateUser(&mocks.MockUserById)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, u)
// 	assert.Equal(t, mocks.MockUser.FullNameEn, u.FullNameEn)
// 	mockRepo.AssertExpectations(t)
// }
