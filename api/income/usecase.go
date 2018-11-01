package income

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

type incomeSum struct {
	Net string
	VAT string
	WHT string
}

func newUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
}

func stringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func floatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func realFloat(f float64) float64 {
	return math.Round(f*100) / 100
}

func calVAT(income string) (string, float64, error) {
	num, err := stringToFloat64(income)
	if err != nil {
		return "", 0.0, err
	}
	vat := num * 0.07
	return floatToString(vat), realFloat(vat), nil
}

func calWHT(income string) (string, float64, error) {
	num, err := stringToFloat64(income)
	if err != nil {
		return "", 0.0, err
	}
	wht := num * 0.03
	return floatToString(wht), realFloat(wht), nil
}

func calIncomeSum(income string, corporateFlag string) (*incomeSum, error) {
	var vat, wht string
	var vatf, whtf float64
	var ins = new(incomeSum)

	total, err := stringToFloat64(income)
	if err != nil {
		return nil, err
	}

	vat, vatf, err = calVAT(income)
	if err != nil {
		return nil, err
	}
	ins.VAT = vat

	if corporateFlag == "Y" {
		wht, whtf, err = calWHT(income)
		if err != nil {
			return nil, err
		}

		net := total + vatf - whtf

		ins.Net = floatToString(net)
		ins.WHT = wht
		return ins, nil
	}

	net := total - vatf
	ins.Net = floatToString(net)
	return ins, nil
}

func getNow() string {
	return time.Now().Format(time.RFC3339)
}

func getCurrentMonth() string {
	y, m, _ := time.Now().Date()
	cm := fmt.Sprintf("%d-%d", y, int(m))
	return cm
}

func (u *usecase) GetIncomeStatusList() ([]*models.IncomeStatus, error) {
	var incomeList []*models.IncomeStatus
	users, err := u.userRepo.GetUser()
	if err != nil {
		return nil, err
	}

	for index, element := range users {
		element.ThaiCitizenID = ""
		incomeUser, err := u.repo.GetIncomeUserNow(element.ID.Hex(), getCurrentMonth())
		income := models.IncomeStatus{User: element}
		incomeList = append(incomeList, &income)
		if err != nil {
			incomeList[index].Status = "N"

		} else {
			incomeList[index].SubmitDate = incomeUser.SubmitDate
			incomeList[index].Status = "Y"
		}

	}
	return incomeList, nil
}

func (u *usecase) GetIncomeByUserIdAndCurrentMonth(userId string) (*models.Income, error) {
	month := getCurrentMonth()
	income, err := u.repo.GetIncomeUserNow(userId, month)
	if err != nil {
		return nil, err
	}
	return income, nil
}
