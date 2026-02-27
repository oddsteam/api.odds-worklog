package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

func NewUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
}

func (u *usecase) GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error) {
	var incomeList []*models.IncomeStatus
	users, err := u.userRepo.GetByRole(role)
	if err != nil {
		return nil, err
	}

	year, month := models.GetYearMonthNow()
	for index, element := range users {
		element.ThaiCitizenID = ""
		element.DailyIncome = ""

		incomeUser, err := u.repo.GetIncomeUserByYearMonth(element.ID.Hex(), year, month)
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

func (u *usecase) GetByRole(role string) ([]*models.User, error) {
	return u.userRepo.GetByRole(role)
}

func (u *usecase) GetUserByID(userId string) (*models.User, error) {
	return u.userRepo.GetByID(userId)
}
