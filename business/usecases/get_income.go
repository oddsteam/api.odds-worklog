package usecases

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type getIncomeUsecase struct {
	repo ForReadingUserIncome
}

func NewGetIncomeUsecase(r ForReadingUserIncome) ForUsingGetIncome {
	return &getIncomeUsecase{r}
}

func (u *getIncomeUsecase) GetIncomeByCurrentMonth(userId string) (*models.Income, error) {
	year, month := models.GetYearMonthNow()
	return u.repo.GetIncomeUserByYearMonth(userId, year, month)
}

func (u *getIncomeUsecase) GetIncomeByAllMonth(userId string) ([]*models.Income, error) {
	listIncome, err := u.repo.GetIncomeByUserIdAllMonth(userId)
	if err != nil {
		return nil, err
	}
	if len(listIncome) == 0 {
		return nil, nil
	}
	for index := range listIncome {
		if listIncome[index].NetSpecialIncome != "" && listIncome[index].NetDailyIncome != "" {
			listIncome[index].NetIncome, err = calIncomeTotal(listIncome[index].NetDailyIncome, listIncome[index].NetSpecialIncome)
			if err != nil {
				return nil, err
			}
		}
	}
	return listIncome, nil
}

func calIncomeTotal(main string, special string) (string, error) {
	ma, err := models.StringToFloat64(main)
	if err != nil {
		return "", err
	}
	sp, err := models.StringToFloat64(special)
	if err != nil {
		return "", err
	}
	return models.FloatToString(ma + sp), nil
}
