package usecases

import (
	"testing"

	siteMock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func claimsFor(u models.User) *models.UserClaims {
	return &models.UserClaims{
		ID:         u.ID.Hex(),
		Role:       u.Role,
		StatusTavi: u.StatusTavi,
	}
}

func TestUsecase_Create(t *testing.T) {
	t.Run("create user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&user).Return(&user, nil)
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, models.ErrInvalidFormat)

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
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
		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		userRes, err := uc.Create(&user)

		assert.EqualError(t, err, models.ErrInvalidFormat.Error())
		assert.Nil(t, userRes)
	})

	t.Run("when user is an exist then create user failed, ErrConflict", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByEmail(user.Email).Return(&user, nil)

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		userRes, err := uc.Create(&user)

		assert.EqualError(t, err, models.ErrConflict.Error())
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

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		u, err := uc.Get()

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, userMock.Users[0].GetFullname(), u[0].GetFullname())
	})

	t.Run("when multiple users share a siteId, then all get Site enriched and SiteID preserved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		siteID := siteMock.MockSite.ID.Hex()
		user1 := userMock.User
		user1.SiteID = siteID
		user2 := userMock.User2
		user2.SiteID = siteID
		users := []*models.User{&user1, &user2}

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockSiteRepo.EXPECT().GetSiteGroup().Return(siteMock.MockSites, nil)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get().Return(users, nil)

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		u, err := uc.Get()

		assert.NoError(t, err)
		assert.Len(t, u, 2)
		for _, got := range u {
			assert.Equal(t, siteID, got.SiteID)
			assert.NotNil(t, got.Site)
			assert.Equal(t, siteMock.MockSite.Name, got.Site.Name)
		}
	})
}

func TestUsecase_GetByRole(t *testing.T) {
	t.Run("when call GetByRole 'corporate', then return list user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
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

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		list, err := uc.GetByRole("individual")

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, userMock.Users, list)
	})
}

func TestUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(userMock.User.ID.Hex()).Return(&userMock.User, nil)

	uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
	u, err := uc.GetByID(userMock.User.ID.Hex())

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, userMock.User.GetFullname(), u.GetFullname())
}

func TestUsecase_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByEmail(userMock.User.Email).Return(&userMock.User, nil)

	uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
	u, err := uc.GetByEmail(userMock.User.Email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, userMock.User.GetEmail(), u.GetEmail())
}

func TestUsecase_GetBySiteID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSiteRepo := siteMock.NewMockRepository(ctrl)
	mockRepo := userMock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetBySiteID("1234567890").Return(userMock.Users, nil)

	uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
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

	uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
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

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		u, err := uc.Update(&userMock.User, claimsFor(userMock.User))

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, userMock.User.GetFullname(), u.GetFullname())
		assert.Equal(t, userMock.User.StartDate, u.StartDate)
	})

	t.Run("when update user invalid role, then retuen erro nil, ErrInvalidUserRole", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(gomock.Any()).Return(&userMock.User, nil)
		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		mu := userMock.User
		mu.Role = "invalid"
		u, err := uc.Update(&mu, claimsFor(userMock.User))

		assert.Nil(t, u)
		assert.EqualError(t, err, models.ErrInvalidUserRole.Error())
	})

	t.Run("when update user invalid vat, then retuen erro nil, ErrInvalidUserVat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(gomock.Any()).Return(&userMock.User, nil)
		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		mu := userMock.User
		mu.Vat = "X"
		u, err := uc.Update(&mu, claimsFor(userMock.User))

		assert.Nil(t, u)
		assert.EqualError(t, err, models.ErrInvalidUserVat.Error())
	})

	t.Run("user-admin can update existing admin user site", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		adminUser := userMock.Admin
		mockRepo.EXPECT().GetByID(adminUser.ID.Hex()).Return(&adminUser, nil)
		mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
			return u, nil
		})

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		req := adminUser
		req.SiteID = "5c0fb860f37e2f8698989cdd"
		u, err := uc.Update(&req, claimsFor(userMock.UserManager))

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, "5c0fb860f37e2f8698989cdd", u.SiteID)
		assert.Equal(t, "admin", u.Role)
	})

	t.Run("user-admin cannot promote user to admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		current := userMock.User
		mockRepo.EXPECT().GetByID(gomock.Any()).Return(&current, nil)
		mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
			return u, nil
		})
		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		req := userMock.User
		req.Role = "admin"
		req.FirstName = "Hacker"
		u, err := uc.Update(&req, claimsFor(userMock.UserManager))

		assert.NoError(t, err)
		assert.Equal(t, "corporate", u.Role)
		assert.Equal(t, userMock.User.FirstName, u.FirstName)
	})

	t.Run("admin can update another user's profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		current := userMock.IndividualUser1
		mockRepo.EXPECT().GetByID(current.ID.Hex()).Return(&current, nil)
		mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
			return u, nil
		})

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		req := current
		req.FirstName = "Updated"
		req.DailyIncome = "9000"
		u, err := uc.Update(&req, claimsFor(userMock.Admin))

		assert.NoError(t, err)
		assert.Equal(t, "Updated", u.FirstName)
		assert.Equal(t, "9000", u.DailyIncome)
	})

	t.Run("user-admin updating another user only applies site", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		current := userMock.IndividualUser1
		current.SiteID = "old-site"
		mockRepo.EXPECT().GetByID(current.ID.Hex()).Return(&current, nil)
		mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
			return u, nil
		})

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		req := current
		req.SiteID = "new-site"
		req.FirstName = "ShouldNotApply"
		req.DailyIncome = "99999"
		u, err := uc.Update(&req, claimsFor(userMock.UserManager))

		assert.NoError(t, err)
		assert.Equal(t, "new-site", u.SiteID)
		assert.Equal(t, userMock.IndividualUser1.FirstName, u.FirstName)
		assert.Equal(t, userMock.IndividualUser1.DailyIncome, u.DailyIncome)
	})

	t.Run("individual can update themselves", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		current := userMock.IndividualUser1
		mockRepo.EXPECT().GetByID(current.ID.Hex()).Return(&current, nil)
		mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(u *models.User) (*models.User, error) {
			return u, nil
		})

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		req := current
		req.Phone = "0812345678"
		u, err := uc.Update(&req, claimsFor(current))

		assert.NoError(t, err)
		assert.Equal(t, "0812345678", u.Phone)
	})

	t.Run("individual cannot update another user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSiteRepo := siteMock.NewMockRepository(ctrl)
		mockRepo := userMock.NewMockRepository(ctrl)
		target := userMock.User
		mockRepo.EXPECT().GetByID(target.ID.Hex()).Return(&target, nil)

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		req := target
		req.FirstName = "Nope"
		u, err := uc.Update(&req, claimsFor(userMock.IndividualUser1))

		assert.Nil(t, u)
		assert.EqualError(t, err, models.ErrPermissionDenied.Error())
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

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
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

		uc := NewManageUsersUsecase(mockRepo, mockSiteRepo)
		mu := userMock.ListUser
		mu[0].User.Role = ""

		u, err := uc.UpdateStatusTavi(mu, mu[0].User.IsAdmin())

		assert.Nil(t, u)
		assert.EqualError(t, err, models.ErrInvalidUserRole.Error())
	})
}
