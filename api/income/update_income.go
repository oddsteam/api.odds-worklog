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

	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome, req.SpecialIncome)
	if err != nil {
		return nil, err
	}
	netIncome, err := utils.StringToFloat64(ins.Net)
	if err != nil {
		return nil, err
	}
	summaryIncome := utils.FloatToString(netIncome)

	income.SubmitDate = time.Now()
	income.TotalIncome = ins.TotalIncome
	income.NetIncome = summaryIncome
	income.VAT = ins.VAT
	income.WHT = ins.WHT
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	u.repo.UpdateIncome(income)

	return income, nil
}
