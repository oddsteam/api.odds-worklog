package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(user.ID.Hex())
	income, err := u.repo.GetIncomeByID(id, user.ID.Hex())
	if err != nil {
		return nil, err
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

	income.SubmitDate = time.Now()
	income.TotalIncome = summaryIncome
	income.NetIncome = netIncomeStr
	income.WorkDate = req.WorkDate
	income.VAT = calNetIncome.VAT
	income.WHT = calNetIncome.WHT
	income.Note = req.Note
	income.SpecialIncome = req.SpecialIncome
	income.NetSpecialIncome = netSpecialIncomeStr
	income.WorkingHours = req.WorkingHours
	u.repo.UpdateIncome(income)

	return income, nil
}
