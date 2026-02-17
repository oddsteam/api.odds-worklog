package usecases

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type addIncomeUsecase struct {
	repo     ForControllingUserIncome
	userRepo ForGettingUserByID
}

func NewAddIncomeUsecase(r ForControllingUserIncome, ur ForGettingUserByID) ForUsingAddIncome {
	return &addIncomeUsecase{r, ur}
}

func (u *addIncomeUsecase) AddIncome(req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	year, month := models.GetYearMonthNow()
	_, err := u.repo.GetIncomeUserByYearMonth(uid, year, month)
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + userDetail.GetName())
	}
	income := models.CreatePayroll(*userDetail, *req, "")
	err = u.repo.AddIncome(income)
	if err != nil {
		return nil, err
	}

	return income, nil
}
