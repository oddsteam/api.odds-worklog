package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestUsecaseUpdateIncome(t *testing.T) {
	t.Run("when update income success it should return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User

		mockUserRepo := mock_usecases.NewMockForGettingUserByID(ctrl)
		mockIncomeRepo := mock_usecases.NewMockForUpdatingUserIncome(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockIncomeRepo.EXPECT().GetIncomeByID(models.MockIncome.ID.Hex(), user.ID.Hex()).Return(&models.MockIncome, nil)
		mockIncomeRepo.EXPECT().UpdateIncome(gomock.Any()).Return(nil)

		uc := NewUpdateIncomeUsecase(mockIncomeRepo, mockUserRepo)
		res, err := uc.UpdateIncome(models.MockIncome.ID.Hex(), &models.MockIncomeReq, user.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncome.UserID, res.UserID)
	})
}
