package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type updateIncomeUsecase struct {
	repo          ForUpdatingUserMonthlyIncome
	userRepo      ForGettingUserByID
	timesheetRepo ForGettingIncomeFromTimesheet
}

func NewUpdateIncomeUsecase(r ForUpdatingUserMonthlyIncome, ur ForGettingUserByID, tr ForGettingIncomeFromTimesheet) ForUsingUpdateIncome {
	return &updateIncomeUsecase{r, ur, tr}
}

func (u *updateIncomeUsecase) UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
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
