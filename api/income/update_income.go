package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	income, err := u.repo.GetIncomeByID(id, uid)
	if err != nil {
		return nil, err
	}

	err = NewIncome(uid).prepareDataForUpdateIncome(*req, *userDetail, income)
	if err != nil {
		return nil, err
	}
	u.repo.UpdateIncome(income)

	return income, nil
}
