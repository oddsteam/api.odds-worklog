package income

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func (u *usecase) AddIncome(req *models.IncomeReq, user *models.User) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(user.ID.Hex())
	year, month := utils.GetYearMonthNow()
	_, err := u.repo.GetIncomeUserByYearMonth(user.ID.Hex(), year, month)
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + userDetail.GetName())
	}
	ins, err := calIncomeSum(req.TotalIncome, userDetail.Vat)
	if err != nil {
		return nil, err
	}

	income := models.Income{
		UserID:      user.ID.Hex(),
		TotalIncome: req.TotalIncome,
		NetIncome:   ins.Net,
		Note:        req.Note,
		VAT:         ins.VAT,
		WHT:         ins.WHT,
	}
	err = u.repo.AddIncome(&income)
	if err != nil {
		return nil, err
	}

	return &income, nil
}
