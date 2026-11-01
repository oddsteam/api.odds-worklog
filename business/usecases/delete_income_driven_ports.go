package usecases

import "time"

// DeleteIncomeByUserYearMonth removes every income document for the user in that month.
// It returns ErrIncomeNotFound when the month has no income document.
type ForDeletingUserMonthlyIncome interface {
	DeleteIncomeByUserYearMonth(userID string, year int, month time.Month) error
}

// DeleteByUserYearMonth removes the income_from_timesheet document for the user in that month.
// It returns ErrIncomeFromTimesheetNotFoundForPeriod when the month has no such document.
type ForDeletingIncomeFromTimesheet interface {
	DeleteByUserYearMonth(userID string, year int, month time.Month) error
}
