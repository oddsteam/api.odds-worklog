package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error) {
	income, err := u.repo.GetIncomeByID(id, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	ins, err := calIncomeSum(req.TotalIncome, user.Vat)
	if err != nil {
		return nil, err
	}

	income.SubmitDate = time.Now()
	income.TotalIncome = req.TotalIncome
	income.NetIncome = ins.Net
	income.VAT = ins.VAT
	income.WHT = ins.WHT
	income.Note = req.Note
	u.repo.UpdateIncome(income)

	return income, nil
}
