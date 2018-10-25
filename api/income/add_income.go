package income

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func (u *usecase) AddIncome(req *models.IncomeReq, user *models.User) (*models.Income, error) {
	userID := user.ID.Hex()

	_, err := u.repo.GetIncomeUserNow(userID, getCurrentMonth())
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + user.FullName)
	}

	ins, err := calIncomeSum(req.TotalIncome, user.CorporateFlag)
	if err != nil {
		return nil, err
	}

	income := models.Income{
		UserID:      userID,
		TotalIncome: req.TotalIncome,
		NetIncome:   ins.Net,
		SubmitDate:  getNow(),
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
