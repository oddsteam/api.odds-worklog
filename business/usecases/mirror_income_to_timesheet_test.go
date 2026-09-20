package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

// The mirrored rows share the income_from_timesheet collection with the ones the timesheet
// consumer writes, so they have to carry the period the same way: it is the key the collection is
// looked up and deduplicated by, and a row without it collides with every other period of the
// same user.
func TestMirrorIncomeToTimesheetStoresThePeriod(t *testing.T) {
	t.Run("stamps the period on a newly mirrored record", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)

		income := models.MockIncome
		repo.EXPECT().GetByUserYearMonth(income.UserID, 2026, time.June).
			Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		repo.EXPECT().Add(gomock.Any()).DoAndReturn(func(rec *models.IncomeFromTimesheet) error {
			assert.Equal(t, 2026, rec.Year)
			assert.Equal(t, 6, rec.Month)
			assert.Equal(t, time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC), rec.SubmitDate.UTC())
			return nil
		})

		err := mirrorIncomeToTimesheet(repo, &income, 2026, time.June)

		assert.NoError(t, err)
	})

	t.Run("stamps the period on an existing mirrored record", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)

		income := models.MockIncome
		// Editing an income from an earlier period must not move the mirrored row into the month the
		// edit happened in — UpdatePayroll has already restamped SubmitDate with now by this point.
		income.SubmitDate = time.Now()
		existing := &models.IncomeFromTimesheet{Income: models.MockIncome, Year: 2026, Month: 6}

		repo.EXPECT().GetByUserYearMonth(income.UserID, 2026, time.June).Return(existing, nil)
		repo.EXPECT().Update(gomock.Any()).DoAndReturn(func(rec *models.IncomeFromTimesheet) error {
			assert.Equal(t, 2026, rec.Year)
			assert.Equal(t, 6, rec.Month)
			assert.Equal(t, time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC), rec.SubmitDate.UTC())
			return nil
		})

		err := mirrorIncomeToTimesheet(repo, &income, 2026, time.June)

		assert.NoError(t, err)
	})
}

// The Income passed in is the one the income collection stores and the one AddIncome/UpdateIncome
// hand back to the caller. Anchoring belongs to the mirrored copy alone — the hand-filled flow
// keeps its own SubmitDate semantics.
func TestMirrorIncomeToTimesheetLeavesTheSourceIncomeAlone(t *testing.T) {
	editedAt := time.Date(2026, 9, 20, 14, 30, 0, 0, time.UTC)

	t.Run("when it creates the mirrored record", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)

		income := models.MockIncome
		income.SubmitDate = editedAt

		repo.EXPECT().GetByUserYearMonth(income.UserID, 2026, time.June).
			Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		repo.EXPECT().Add(gomock.Any()).Return(nil)

		err := mirrorIncomeToTimesheet(repo, &income, 2026, time.June)

		assert.NoError(t, err)
		assert.Equal(t, editedAt, income.SubmitDate)
	})

	t.Run("when it updates the mirrored record", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)

		income := models.MockIncome
		income.SubmitDate = editedAt
		existing := &models.IncomeFromTimesheet{Income: models.MockIncome, Year: 2026, Month: 6}

		repo.EXPECT().GetByUserYearMonth(income.UserID, 2026, time.June).Return(existing, nil)
		repo.EXPECT().Update(gomock.Any()).Return(nil)

		err := mirrorIncomeToTimesheet(repo, &income, 2026, time.June)

		assert.NoError(t, err)
		assert.Equal(t, editedAt, income.SubmitDate)
	})
}
