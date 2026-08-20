package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestIncomeFromTimesheetSource(t *testing.T) {
	startDate := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)

	t.Run("unwraps the embedded Income of each timesheet record", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		first := &models.IncomeFromTimesheet{Income: models.MockIncome}
		first.UserID = "user-1"
		second := &models.IncomeFromTimesheet{Income: models.MockIncome}
		second.UserID = "user-2"

		repo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return([]*models.IncomeFromTimesheet{first, second}, nil)

		source := NewIncomeFromTimesheetSource(repo)
		incomes, err := source.GetAllIncomeByRoleStartDateAndEndDate("individual", startDate, endDate)

		assert.NoError(t, err)
		assert.Len(t, incomes, 2)
		assert.Equal(t, "user-1", incomes[0].UserID)
		assert.Equal(t, "user-2", incomes[1].UserID)
	})

	t.Run("skips nil records rather than dereferencing them", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		record := &models.IncomeFromTimesheet{Income: models.MockIncome}
		record.UserID = "user-1"

		repo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return([]*models.IncomeFromTimesheet{nil, record}, nil)

		source := NewIncomeFromTimesheetSource(repo)
		incomes, err := source.GetAllIncomeByRoleStartDateAndEndDate("individual", startDate, endDate)

		assert.NoError(t, err)
		assert.Len(t, incomes, 1)
		assert.Equal(t, "user-1", incomes[0].UserID)
	})

	t.Run("propagates an error from the repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		repo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return(nil, assert.AnError)

		source := NewIncomeFromTimesheetSource(repo)
		_, err := source.GetAllIncomeByRoleStartDateAndEndDate("individual", startDate, endDate)

		assert.ErrorIs(t, err, assert.AnError)
	})
}
