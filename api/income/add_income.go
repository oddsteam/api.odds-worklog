package income

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func (u *usecase) AddIncome(req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	year, month := utils.GetYearMonthNow()
	_, err := u.repo.GetIncomeUserByYearMonth(uid, year, month)
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + userDetail.GetName())
	}
	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome,userDetail.GetRole())
	if err != nil {
		return nil, err
	}
	insSpecial, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome,userDetail.GetRole())
	if err != nil {
		return nil, err
	}
	summaryIncome, err := calSummary(ins.TotalIncome, insSpecial.TotalIncome)
	if err != nil {
		return nil, err
	}
	summaryWht, err := calSummary(ins.WHT, insSpecial.WHT)
	if err != nil {
		return nil, err
	}
	var summaryVat string
	if userDetail.Vat != "N" {
		summaryVat, err = calSummary(ins.VAT, insSpecial.VAT)
		if err != nil {
			return nil, err
		}
	} else {
		summaryVat = ""
	}

	income := models.Income{
		UserID:           uid,
		TotalIncome:      summaryIncome,
		NetIncome:        ins.Net,
		NetSpecialIncome: insSpecial.Net,
		Note:             req.Note,
		VAT:              summaryVat,
		WHT:              summaryWht,
		WorkDate:         req.WorkDate,
		SpecialIncome:    req.SpecialIncome,
		WorkingHours:     req.WorkingHours,
	}
	err = u.repo.AddIncome(&income)
	if err != nil {
		return nil, err
	}

	return &income, nil
}
