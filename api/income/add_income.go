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
	calNetIncome, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome)
	if err != nil {
		return nil, err
	}
	netIncome, err := utils.StringToFloat64(calNetIncome.Net)
	if err != nil {
		return nil, err
	}
	calNetSpecialIncome, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome)
	if err != nil {
		return nil, err
	}
	netSpecialIncome, err := utils.StringToFloat64(calNetSpecialIncome.Net)
	if err != nil {
		return nil, err
	}

	netIncomeStr := utils.FloatToString(netIncome)
	netSpecialIncomeStr := utils.FloatToString(netSpecialIncome)
	summaryIncome := utils.FloatToString(netIncome + netSpecialIncome)

	income := models.Income{
		UserID:           user.ID.Hex(),
		TotalIncome:      summaryIncome,
		NetIncome:        netIncomeStr,
		Note:             req.Note,
		VAT:              calNetIncome.VAT,
		WHT:              calNetIncome.WHT,
		WorkDate:         req.WorkDate,
		SpecialIncome:    req.SpecialIncome,
		NetSpecialIncome: netSpecialIncomeStr,
		WorkingHours:     req.WorkingHours,
	}
	err = u.repo.AddIncome(&income)
	if err != nil {
		return nil, err
	}

	return &income, nil
}
