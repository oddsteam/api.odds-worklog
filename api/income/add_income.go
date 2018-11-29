package income

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func (u *usecase) AddIncome(req *models.IncomeReq, user *models.User) (*models.Income, error) {
	userID := user.ID.Hex()
	year, month := utils.GetYearMonthNow()
	_, err := u.repo.GetIncomeUserByYearMonth(userID, year, month)
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + user.FullNameEn)
	}

	ins, err := calIncomeSum(req.TotalIncome, user.Vat)
	if err != nil {
		return nil, err
	}

	income := models.Income{
		UserID:      userID,
		TotalIncome: req.TotalIncome,
		NetIncome:   ins.Net,
		SubmitDate:  time.Now(),
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
