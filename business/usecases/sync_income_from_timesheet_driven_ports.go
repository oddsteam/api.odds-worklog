package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// GetByUserYearMonth must return ErrIncomeFromTimesheetNotFoundForPeriod (not a raw driver
// error) when no record exists yet for the given user+year+month — any other non-nil error
// is a real failure and must be propagated, not treated as "not found."
type ForGettingIncomeFromTimesheet interface {
	GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeFromTimesheet, error)
	Add(income *models.IncomeFromTimesheet) error
	Update(income *models.IncomeFromTimesheet) error
}

// GetByEmail must return ErrTimesheetUserNotFound (not a raw driver error) when no user
// matches the email — any other non-nil error is a real failure and must be propagated.
//
// DisableIncomeEntry sets only the user's disableIncomeEntry flag, leaving the rest of the
// document as it is, so it never overwrites a profile edit made in between.
type ForGettingTimesheetUser interface {
	GetByEmail(email string) (*models.User, error)
	DisableIncomeEntry(userID string) error
}

// Save persists the raw event payload unconditionally, before any other processing.
type ForLoggingTimesheetEvent interface {
	Save(evt models.TimesheetMonthlySummaryEvent) error
}
