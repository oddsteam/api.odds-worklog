package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestListIncomeStatusUsecase_GetIncomeStatusList(t *testing.T) {
	t.Run("when get list income success it should be return list income where status is Y", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockIncomeRepo := mock_usecases.NewMockForReadingUserIncome(ctrl)
		year, month := models.GetYearMonthNow()
		mockIncomeRepo.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month).Return(&models.MockIncome, nil)
		mockIncomeRepo.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month).Return(&models.MockIncome, nil)

		mockUserRepo := mock_usecases.NewMockForGettingUsersByRole(ctrl)
		mockUserRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewListIncomeStatusUsecase(mockIncomeRepo, mockUserRepo)
		res, err := uc.GetIncomeStatusList("corporate", false)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncomeStatusList[0].Status, res[0].Status)
	})
}
