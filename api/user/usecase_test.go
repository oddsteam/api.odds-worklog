package user

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	mock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_CreateUser(t *testing.T) {
	t.Run("create user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.MockUser

		mockRepo := mock.NewMockRepository(ctrl)
		mockRepo.EXPECT().CreateUser(&user).Return(&user, nil)
		mockRepo.EXPECT().GetUserByEmail(user.Email).Return(nil, utils.ErrNotFound)

		uc := NewUsecase(mockRepo)
		userRes, err := uc.CreateUser(&user)

		assert.NoError(t, err)
		assert.NotNil(t, userRes)
		assert.Equal(t, user.ID, userRes.ID)
	})

	t.Run("when email is invalid then create user failed, ErrInvalidFormat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.MockUser
		user.Email = "abc"

		mockRepo := mock.NewMockRepository(ctrl)
		uc := NewUsecase(mockRepo)
		userRes, err := uc.CreateUser(&user)

		assert.EqualError(t, err, utils.ErrInvalidFormat.Error())
		assert.Nil(t, userRes)
	})

	t.Run("when user is an exist then create user failed, ErrConflict", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.MockUser
		mockRepo := mock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetUserByEmail(user.Email).Return(&user, nil)

		uc := NewUsecase(mockRepo)
		userRes, err := uc.CreateUser(&user)

		assert.EqualError(t, err, utils.ErrConflict.Error())
		assert.NotNil(t, userRes)
	})

}

func TestUsecase_GetUser(t *testing.T) {
	t.Run("when call GetUser, then user not nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetUser().Return(mock.MockUsers, nil)

		uc := NewUsecase(mockRepo)
		u, err := uc.GetUser()

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, mock.MockUsers[0].GetFullname(), u[0].GetFullname())
	})
}

func TestUsecase_GetUserByRole(t *testing.T) {
	t.Run("when call GetUserByRole 'corporate', then return list user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetUserByRole("corporate").Return(mock.MockUsers, nil)

		uc := NewUsecase(mockRepo)
		list, err := uc.GetUserByRole("corporate")

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, mock.MockUsers, list)
	})

	t.Run("when call GetUserByRole 'individual', then return list user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetUserByRole("individual").Return(mock.MockUsers, nil)

		uc := NewUsecase(mockRepo)
		list, err := uc.GetUserByRole("individual")

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, mock.MockUsers, list)
	})
}

func TestUsecase_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetUserByID("1234567890").Return(&mock.MockUserById, nil)

	uc := NewUsecase(mockRepo)
	u, err := uc.GetUserByID(string(mock.MockUserById.ID))

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mock.MockUserById.GetFullname(), u.GetFullname())
}

func TestUsecase_GetUserBySiteID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetUserBySiteID("1234567890").Return(mock.MockUsers, nil)

	uc := NewUsecase(mockRepo)
	users, err := uc.GetUserBySiteID("1234567890")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, mock.MockUsers, users)
}

func TestUsecase_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().DeleteUser("1234567890").Return(nil)

	uc := NewUsecase(mockRepo)
	u := uc.DeleteUser(string(mock.MockUserById.ID))

	assert.Equal(t, nil, u)
}

func TestUsecase_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(&mock.MockUserById, nil)

	uc := NewUsecase(mockRepo)
	u, err := uc.UpdateUser(&mock.MockUserById, nil)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mock.MockUser.GetFullname(), u.GetFullname())
}

func TestUsecase_getTranscriptFilename(t *testing.T) {
	u := mock.MockUser

	filename := getTranscriptFilename(&u)
	assert.NotEmpty(t, filename)

	path := "files/transcripts"
	filenameExp := fmt.Sprintf("%s/transcript_%s_%s_", path, strings.ToUpper(u.FirstName), strings.ToUpper(u.LastName))
	assert.Contains(t, filename, filenameExp)
	assert.Contains(t, filename, ".pdf")
	assert.Equal(t, len(filenameExp)+12, len(filename))
}
