package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type updateIncomeUsecase struct {
	repo     ForUpdatingUserIncome
	userRepo ForGettingUserByID
}

func NewUpdateIncomeUsecase(r ForUpdatingUserIncome, ur ForGettingUserByID) ForUsingUpdateIncome {
	return &updateIncomeUsecase{r, ur}
}

func (u *updateIncomeUsecase) UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	income, err := u.repo.GetIncomeByID(id, uid)
	if err != nil {
		return nil, err
	}
	income = models.UpdatePayroll(*userDetail, *req, "", income)
	u.repo.UpdateIncome(income)
	return income, nil
}
