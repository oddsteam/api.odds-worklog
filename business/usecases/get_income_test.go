package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/business/models/mock"
)

func TestUsecaseGetIncomeByCurrentMonth(t *testing.T) {
	t.Run("when get income by user id current month success it should return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		year, month := models.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(models.MockIncome.UserID, year, month).Return(&models.MockIncome, nil)

		uc := NewGetIncomeUsecase(mockRepoIncome)
		res, err := uc.GetIncomeByCurrentMonth(models.MockIncome.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncome.SubmitDate, res.SubmitDate)
	})
}

func TestUsecaseGetIncomeByAllMonth(t *testing.T) {
	t.Run("when get income by user id all month success it should return income list with calculated net income", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserIdAllMonth(models.MockIncome.UserID).Return(models.MockIncomeList, nil)

		uc := NewGetIncomeUsecase(mockRepoIncome)
		res, err := uc.GetIncomeByAllMonth(models.MockIncome.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncome.SubmitDate, res[0].SubmitDate)
		assert.Equal(t, "50440.00", res[1].NetIncome)
	})
}

func TestUsecaseGetIncomeByAllMonthCaseNoNetSpecialIncome(t *testing.T) {
	t.Run("when get income by user id all month and no net special income it should return income list without recalculation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserIdAllMonth(models.MockIncome.UserID).Return(models.MockIncomeListNoNetSpecialIncome, nil)

		uc := NewGetIncomeUsecase(mockRepoIncome)
		res, err := uc.GetIncomeByAllMonth(models.MockIncome.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "00.00", res[0].NetIncome)
		assert.Equal(t, "50440.00", res[1].NetIncome)
	})
}
