package usecases

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// ErrIncomeNotFound is returned when an id has no record in the income collection for this user.
// It is its own error so the api layer can answer 404 instead of 500 — the case shows up when a
// caller edits with an id read from another collection, e.g. the income_from_timesheet mirror,
// whose ids never exist in income.
var ErrIncomeNotFound = errors.New("income: no record for this id and user")

type updateIncomeUsecase struct {
	repo          ForUpdatingUserMonthlyIncome
	userRepo      ForGettingUserByID
	timesheetRepo ForGettingIncomeFromTimesheet
	siteRepo      ForGettingSiteByID
}

func NewUpdateIncomeUsecase(r ForUpdatingUserMonthlyIncome, ur ForGettingUserByID, tr ForGettingIncomeFromTimesheet, sr ForGettingSiteByID) ForUsingUpdateIncome {
	return &updateIncomeUsecase{r, ur, tr, sr}
}

func (u *updateIncomeUsecase) UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	attachSite(userDetail, u.siteRepo)
	income, err := u.repo.GetIncomeByID(id, uid)
	if err != nil {
		return nil, err
	}
	// The path id belongs to the income collection, so the mirrored record has to be found by
	// period instead — read off before UpdatePayroll, which restamps SubmitDate with now.
	submitted := income.SubmitDate.UTC()
	income = models.UpdatePayroll(*userDetail, *req, req.Note, income)
	if err := u.repo.UpdateIncome(income); err != nil {
		return nil, err
	}
	if err := mirrorIncomeToTimesheet(u.timesheetRepo, income, submitted.Year(), submitted.Month()); err != nil {
		return nil, err
	}
	return income, nil
}
