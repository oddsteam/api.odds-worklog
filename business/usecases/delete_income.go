package usecases

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var (
	ErrUserNotFound            = errors.New("user: not found")
	ErrIncomeUserNotIndividual = errors.New("income: user is not an individual")
)

type deleteIncomeUsecase struct {
	users     ForGettingUserByID
	income    ForDeletingUserMonthlyIncome
	timesheet ForDeletingIncomeFromTimesheet
}

func NewDeleteIncomeUsecase(users ForGettingUserByID, income ForDeletingUserMonthlyIncome, timesheet ForDeletingIncomeFromTimesheet) ForUsingDeleteIncome {
	return &deleteIncomeUsecase{users, income, timesheet}
}

// DeleteCurrentMonthIncome removes this month's income for an individual user from both
// the income collection and the income_from_timesheet mirror, so either list shows the
// user as not submitted. A later timesheet event can write the timesheet copy again.
func (u *deleteIncomeUsecase) DeleteCurrentMonthIncome(userID string) error {
	user, err := u.users.GetByID(userID)
	if err != nil {
		return err
	}
	if user.GetRole() != "individual" {
		return ErrIncomeUserNotIndividual
	}

	year, month := models.GetYearMonthNow()
	incomeErr := u.income.DeleteIncomeByUserYearMonth(userID, year, month)
	timesheetErr := u.timesheet.DeleteByUserYearMonth(userID, year, month)

	incomeMissing := errors.Is(incomeErr, ErrIncomeNotFound)
	timesheetMissing := errors.Is(timesheetErr, ErrIncomeFromTimesheetNotFoundForPeriod)
	if incomeMissing && timesheetMissing {
		return ErrIncomeNotFound
	}
	if incomeErr != nil && !incomeMissing {
		return incomeErr
	}
	if timesheetErr != nil && !timesheetMissing {
		return timesheetErr
	}
	return nil
}
