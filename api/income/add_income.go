package income

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func (u *usecase) AddIncome(req *models.IncomeReq, user *models.User) (*models.Income, error) {
	userID := user.ID.Hex()

	_, err := u.repo.GetIncomeUserNow(userID, utils.GetCurrentMonth())
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + user.FullNameEn)
	}

	ins, err := calIncomeSum(req.TotalIncome, user.CorporateFlag)
	if err != nil {
		return nil, err
	}

	income := models.Income{
		UserID:      userID,
		TotalIncome: req.TotalIncome,
		NetIncome:   ins.Net,
		SubmitDate:  utils.GetNow(),
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
