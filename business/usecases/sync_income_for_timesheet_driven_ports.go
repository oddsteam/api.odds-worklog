package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// GetByUserYearMonth must return ErrIncomeForTimesheetNotFoundForPeriod (not a raw driver
// error) when no record exists yet for the given user+year+month — any other non-nil error
// is a real failure and must be propagated, not treated as "not found."
type ForGettingIncomeForTimesheet interface {
	GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeForTimesheet, error)
	Add(income *models.IncomeForTimesheet) error
	Update(income *models.IncomeForTimesheet) error
}

// GetByEmail must return ErrTimesheetUserNotFound (not a raw driver error) when no user
// matches the email — any other non-nil error is a real failure and must be propagated.
type ForGettingTimesheetUser interface {
	GetByEmail(email string) (*models.User, error)
}
