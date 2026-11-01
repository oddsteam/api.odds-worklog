package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type stubUserLookup struct {
	user *models.User
	err  error
}

func (s *stubUserLookup) GetByID(id string) (*models.User, error) {
	return s.user, s.err
}

type stubIncomeDeleter struct {
	err    error
	userID string
	year   int
	month  time.Month
	called bool
}

func (s *stubIncomeDeleter) DeleteIncomeByUserYearMonth(userID string, year int, month time.Month) error {
	s.called = true
	s.userID = userID
	s.year = year
	s.month = month
	return s.err
}

type stubTimesheetDeleter struct {
	err    error
	userID string
	year   int
	month  time.Month
	called bool
}

func (s *stubTimesheetDeleter) DeleteByUserYearMonth(userID string, year int, month time.Month) error {
	s.called = true
	s.userID = userID
	s.year = year
	s.month = month
	return s.err
}

func TestDeleteCurrentMonthIncome(t *testing.T) {
	userID := "5bbcf2f90fd2df527bc39531"
	individual := &models.User{Role: "individual"}
	year, month := models.GetYearMonthNow()

	t.Run("deletes this month from income and income_from_timesheet", func(t *testing.T) {
		income := &stubIncomeDeleter{}
		timesheet := &stubTimesheetDeleter{}
		uc := NewDeleteIncomeUsecase(&stubUserLookup{user: individual}, income, timesheet)

		err := uc.DeleteCurrentMonthIncome(userID)

		assert.NoError(t, err)
		assert.Equal(t, userID, income.userID)
		assert.Equal(t, year, income.year)
		assert.Equal(t, month, income.month)
		assert.Equal(t, userID, timesheet.userID)
		assert.Equal(t, year, timesheet.year)
		assert.Equal(t, month, timesheet.month)
	})

	t.Run("succeeds when only the income record exists", func(t *testing.T) {
		uc := NewDeleteIncomeUsecase(
			&stubUserLookup{user: individual},
			&stubIncomeDeleter{},
			&stubTimesheetDeleter{err: ErrIncomeFromTimesheetNotFoundForPeriod},
		)

		assert.NoError(t, uc.DeleteCurrentMonthIncome(userID))
	})

	t.Run("succeeds when only the timesheet record exists", func(t *testing.T) {
		uc := NewDeleteIncomeUsecase(
			&stubUserLookup{user: individual},
			&stubIncomeDeleter{err: ErrIncomeNotFound},
			&stubTimesheetDeleter{},
		)

		assert.NoError(t, uc.DeleteCurrentMonthIncome(userID))
	})

	t.Run("returns not found when neither record exists", func(t *testing.T) {
		uc := NewDeleteIncomeUsecase(
			&stubUserLookup{user: individual},
			&stubIncomeDeleter{err: ErrIncomeNotFound},
			&stubTimesheetDeleter{err: ErrIncomeFromTimesheetNotFoundForPeriod},
		)

		err := uc.DeleteCurrentMonthIncome(userID)
		assert.ErrorIs(t, err, ErrIncomeNotFound)
	})

	t.Run("refuses a user who is not individual", func(t *testing.T) {
		income := &stubIncomeDeleter{}
		uc := NewDeleteIncomeUsecase(
			&stubUserLookup{user: &models.User{Role: "corporate"}},
			income,
			&stubTimesheetDeleter{},
		)

		err := uc.DeleteCurrentMonthIncome(userID)
		assert.ErrorIs(t, err, ErrIncomeUserNotIndividual)
		assert.False(t, income.called)
	})

	t.Run("returns the lookup error when the user does not exist", func(t *testing.T) {
		uc := NewDeleteIncomeUsecase(
			&stubUserLookup{err: ErrUserNotFound},
			&stubIncomeDeleter{},
			&stubTimesheetDeleter{},
		)

		err := uc.DeleteCurrentMonthIncome(userID)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("returns a repository failure", func(t *testing.T) {
		repoErr := errors.New("db down")
		uc := NewDeleteIncomeUsecase(
			&stubUserLookup{user: individual},
			&stubIncomeDeleter{err: repoErr},
			&stubTimesheetDeleter{},
		)

		err := uc.DeleteCurrentMonthIncome(userID)
		assert.ErrorIs(t, err, repoErr)
	})
}
