package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(user.ID.Hex())
	income, err := u.repo.GetIncomeByID(id, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome)
	if err != nil {
		return nil, err
	}
	insSpecial, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome)
	if err != nil {
		return nil, err
	}
	summaryIncome, err := calTotalIncome(ins.TotalIncome, insSpecial.TotalIncome)
	if err != nil {
		return nil, err
	}
	summaryWht, err := calSummaryWht(ins.WHT, insSpecial.WHT)
	if err != nil {
		return nil, err
	}
	var summaryVat string
	if userDetail.Vat != "N" {
		summaryVat, err = calSummaryVat(ins.VAT, insSpecial.VAT)
		if err != nil {
			return nil, err
		}
	} else {
		summaryVat = ""
	}

	income.SubmitDate = time.Now()
	income.TotalIncome = summaryIncome
	income.NetIncome = ins.Net
	income.NetSpecialIncome = insSpecial.Net
	income.VAT = summaryVat
	income.WHT = summaryWht
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours
	u.repo.UpdateIncome(income)
	return income, nil
}
