package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestIncomeFromTimesheetUserSource(t *testing.T) {
	t.Run("GetIncomeUserByYearMonth unwraps the embedded Income", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForReadingIncomeFromTimesheetByUser(ctrl)

		record := &models.IncomeFromTimesheet{Income: models.MockIncome}
		repo.EXPECT().GetByUserYearMonth("user-1", 2026, time.June).Return(record, nil)

		source := NewIncomeFromTimesheetUserSource(repo)
		income, err := source.GetIncomeUserByYearMonth("user-1", 2026, time.June)

		assert.NoError(t, err)
		assert.Equal(t, models.MockIncome.SubmitDate, income.SubmitDate)
	})

	t.Run("GetIncomeUserByYearMonth propagates an error from the repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForReadingIncomeFromTimesheetByUser(ctrl)

		repo.EXPECT().GetByUserYearMonth("user-1", 2026, time.June).Return(nil, assert.AnError)

		source := NewIncomeFromTimesheetUserSource(repo)
		_, err := source.GetIncomeUserByYearMonth("user-1", 2026, time.June)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("GetIncomeByUserIdAllMonth unwraps the embedded Income of each record and skips nils", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForReadingIncomeFromTimesheetByUser(ctrl)

		first := &models.IncomeFromTimesheet{Income: models.MockIncome}
		first.UserID = "user-1"

		repo.EXPECT().GetByUserIdAllMonth("user-1").Return([]*models.IncomeFromTimesheet{nil, first}, nil)

		source := NewIncomeFromTimesheetUserSource(repo)
		incomes, err := source.GetIncomeByUserIdAllMonth("user-1")

		assert.NoError(t, err)
		assert.Len(t, incomes, 1)
		assert.Equal(t, "user-1", incomes[0].UserID)
	})

	t.Run("GetIncomeByUserIdAllMonth propagates an error from the repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForReadingIncomeFromTimesheetByUser(ctrl)

		repo.EXPECT().GetByUserIdAllMonth("user-1").Return(nil, assert.AnError)

		source := NewIncomeFromTimesheetUserSource(repo)
		_, err := source.GetIncomeByUserIdAllMonth("user-1")

		assert.ErrorIs(t, err, assert.AnError)
	})
}
