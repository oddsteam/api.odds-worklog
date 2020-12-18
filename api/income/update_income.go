package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	income, err := u.repo.GetIncomeByID(id, uid)
	if err != nil {
		return nil, err
	}

	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome, userDetail.GetRole())
	if err != nil {
		return nil, err
	}
	insSpecial, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome, userDetail.GetRole())
	if err != nil {
		return nil, err
	}
	summaryIncome, err := calSummary(ins.TotalIncome, insSpecial.TotalIncome)
	if err != nil {
		return nil, err
	}
	summaryNetIncome, err := calSummary(ins.Net, insSpecial.Net)
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

	income.SubmitDate = time.Now()
	income.TotalIncome = summaryIncome
	income.NetIncome = summaryNetIncome
	income.NetSpecialIncome = insSpecial.Net
	income.NetDailyIncome = ins.Net
	income.VAT = summaryVat
	income.WHT = summaryWht
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours
	u.repo.UpdateIncome(income)
	return income, nil
}
