package user

import (
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	siteMock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Create(t *testing.T) {
	t.Run("create user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&user).Return(&user, nil)
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, utils.ErrNotFound)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		userRes, err := uc.Create(&user)

		assert.NoError(t, err)
		assert.NotNil(t, userRes)
		assert.Equal(t, user.ID, userRes.ID)
	})

	t.Run("when email is invalid then create user failed, ErrInvalidFormat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		user.Email = "abc"

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		uc := NewUsecase(mockRepo, mockSiteRepo)
		userRes, err := uc.Create(&user)

		assert.EqualError(t, err, utils.ErrInvalidFormat.Error())
		assert.Nil(t, userRes)
	})

	t.Run("when user is an exist then create user failed, ErrConflict", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByEmail(user.Email).Return(&user, nil)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		userRes, err := uc.Create(&user)

		assert.EqualError(t, err, utils.ErrConflict.Error())
		assert.NotNil(t, userRes)
	})

}

func TestUsecase_Get(t *testing.T) {
	t.Run("when call Get, then user not nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockSiteRepo.EXPECT().GetSiteGroup().Return(siteMock.MockSites, nil)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get().Return(userMock.Users, nil)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		u, err := uc.Get()

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, userMock.Users[0].GetFullname(), u[0].GetFullname())
	})
}

func TestUsecase_GetByRole(t *testing.T) {
	t.Run("when call GetByRole 'corporate', then return list user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		list, err := uc.GetByRole("corporate")

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, userMock.Users, list)
	})

	t.Run("when call GetByRole 'individual', then return list user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByRole("individual").Return(userMock.Users, nil)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		list, err := uc.GetByRole("individual")

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, userMock.Users, list)
	})
}

func TesTUsercase_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByEmail(userMock.User.Email).Return(&userMock.User, nil)

	uc := NewUsecase(mockRepo, mockSiteRepo)
	u, err := uc.GetByEmail(userMock.User.Email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, userMock.User.GetEmail(), u.GetEmail())
}

func TestUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(userMock.User.ID.Hex()).Return(&userMock.User, nil)

	uc := NewUsecase(mockRepo, mockSiteRepo)
	u, err := uc.GetByID(userMock.User.ID.Hex())

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, userMock.User.GetFullname(), u.GetFullname())
}
func TestUsecase_GetBySiteID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetBySiteID("1234567890").Return(userMock.Users, nil)

	uc := NewUsecase(mockRepo, mockSiteRepo)
	users, err := uc.GetBySiteID("1234567890")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, userMock.Users, users)
}

func TestUsecase_Delete_Should_Move_To_Archived_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().Delete(userMock.User.ID.Hex()).Return(nil)
	mockRepo.EXPECT().CreateArchivedUser(userMock.User).Return(nil, nil)
	mockRepo.EXPECT().GetByID(userMock.User.ID.Hex()).Return(&userMock.User, nil)

	uc := NewUsecase(mockRepo, mockSiteRepo)
	u := uc.Delete(userMock.User.ID.Hex())

	assert.Equal(t, nil, u)
}

func TestUsecase_Update(t *testing.T) {
	t.Run("update user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(gomock.Any()).Return(&userMock.User, nil)
		mockRepo.EXPECT().Update(gomock.Any()).Return(&userMock.User, nil)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		u, err := uc.Update(&userMock.User, userMock.User.IsAdmin())

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, userMock.User.GetFullname(), u.GetFullname())
		assert.Equal(t, userMock.User.StartDate, u.StartDate)
	})

	t.Run("Bank account number with - and special char will create a bad batch file for bank system. This will fail the batch transfer process in the bank, causing the delay for all members to receive income. Therefore, we will remove - and special char from the bank account number!", func(t *testing.T) {
		userFromRequest := models.User{
			BankAccountNumber: "้1234-123-999",
		}

		user := NewUser(userMock.User)
		err := user.prepareDataForUpdateFrom(userFromRequest)

		assert.NoError(t, err)
		assert.Equal(t, user.data.BankAccountNumber, "1234123999")
	})

	t.Run("when update user invalid role, then retuen erro nil, ErrInvalidUserRole", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		uc := NewUsecase(mockRepo, mockSiteRepo)
		mu := userMock.User
		mu.Role = ""
		u, err := uc.Update(&mu, mu.IsAdmin())

		assert.Nil(t, u)
		assert.EqualError(t, err, utils.ErrInvalidUserRole.Error())
	})

	t.Run("when update user invalid vat, then retuen erro nil, ErrInvalidUserVat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		uc := NewUsecase(mockRepo, mockSiteRepo)
		mu := userMock.User
		mu.Vat = ""
		u, err := uc.Update(&mu, mu.IsAdmin())

		assert.Nil(t, u)
		assert.EqualError(t, err, utils.ErrInvalidUserVat.Error())
	})
}

func TestUsecase_UpdateStatusTavi(t *testing.T) {
	t.Run("update status tavi success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(gomock.Any()).Return(&userMock.User, nil)
		mockRepo.EXPECT().Update(gomock.Any()).Return(&userMock.User, nil)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		u, err := uc.UpdateStatusTavi(userMock.ListUser, userMock.User.IsAdmin())

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, userMock.User.GetFullname(), u[0].GetFullname())
	})

	t.Run("when update status tavi invalid role, then retuen erro nil, ErrInvalidUserRole", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepo, mockSiteRepo)
		mu := userMock.ListUser
		mu[0].User.Role = ""

		u, err := uc.UpdateStatusTavi(mu, mu[0].User.IsAdmin())

		assert.Nil(t, u)
		assert.EqualError(t, err, utils.ErrInvalidUserRole.Error())
	})

}
