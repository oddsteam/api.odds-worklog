package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error) {
	income, err := u.repo.GetIncomeByID(id, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	ins, err := calIncomeSum(req.TotalIncome, user.CorporateFlag)
	if err != nil {
		return nil, err
	}

	income.SubmitDate = utils.GetNow()
	income.TotalIncome = req.TotalIncome
	income.NetIncome = ins.Net
	income.VAT = ins.VAT
	income.WHT = ins.WHT
	income.Note = req.Note
	u.repo.UpdateIncome(income)

	return income, nil
}
