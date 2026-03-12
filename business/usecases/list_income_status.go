package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type listIncomeStatusUsecase struct {
	incomeRepo ForReadingUserIncome
	userRepo   ForListingUsersByRole
}

func NewListIncomeStatusUsecase(incomeRepo ForReadingUserIncome, userRepo ForListingUsersByRole) ForUsingListIncomeStatus {
	return &listIncomeStatusUsecase{incomeRepo: incomeRepo, userRepo: userRepo}
}

func (u *listIncomeStatusUsecase) GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error) {
	var incomeList []*models.IncomeStatus
	users, err := u.userRepo.GetByRole(role)
	if err != nil {
		return nil, err
	}

	year, month := models.GetYearMonthNow()
	for index, element := range users {
		element.ThaiCitizenID = ""
		element.DailyIncome = ""

		incomeUser, err := u.incomeRepo.GetIncomeUserByYearMonth(element.ID.Hex(), year, month)
		income := models.IncomeStatus{User: element}
		incomeList = append(incomeList, &income)
		if !isAdmin {
			element.ID = ""
		}
		if err != nil {
			incomeList[index].Status = "N"
		} else {
			incomeList[index].WorkDate = incomeUser.WorkDate
			incomeList[index].WorkingHours = incomeUser.WorkingHours
			incomeList[index].SubmitDate = incomeUser.SubmitDate.Format(time.RFC3339)
			incomeList[index].Status = "Y"
		}
	}
	return incomeList, nil
}
