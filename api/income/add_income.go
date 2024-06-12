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
	i := NewIncome(uid)
	i.prepareDataForAddIncome(*req, *userDetail)

	income := models.Income{
		UserID:           uid,
		TotalIncome:      i.TotalIncomeStr,
		NetIncome:        i.NetIncomeStr,
		NetSpecialIncome: i.NetSpecialIncomeStr,
		NetDailyIncome:   i.NetDailyIncomeStr,
		Note:             req.Note,
		VAT:              i.VATStr,
		WHT:              i.WHTStr,
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
