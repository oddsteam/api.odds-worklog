package income

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/business/models/mock"
)

func TestUsecaseGetListIncome(t *testing.T) {
	t.Run("when get list income success it should be return list income where status is Y", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		year, month := models.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month).Return(&models.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month).Return(&models.MockIncome, nil)

		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeStatusList("corporate", false)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncomeStatusList[0].Status, res[0].Status)

	})
}
