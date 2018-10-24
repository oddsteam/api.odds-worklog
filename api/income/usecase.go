package income

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"

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

func (u *usecase) AddIncome(req *models.IncomeReq, id string) (*models.IncomeRes, error) {
	us, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	y, m, _ := t.Date()
	cm := fmt.Sprintf("%d-%d", y, int(m))

	_, err = u.repo.GetIncomeUserNow(us.ID.Hex(), cm)
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + us.FullName)
	}

	ins, err := calIncomeSum(req.TotalIncome, us.CorporateFlag)
	if err != nil {
		return nil, err
	}

	income := models.Income{
		UserID:      id,
		TotalIncome: req.TotalIncome,
		NetIncome:   ins.Net,
		SubmitDate:  t.Format(time.RFC3339),
		Note:        req.Note,
		VAT:         ins.VAT,
		WHT:         ins.WHT,
	}
	err = u.repo.AddIncome(&income)
	if err != nil {
		return nil, err
	}

	res := &models.IncomeRes{
		Income: &income,
		Status: "Y",
	}
	return res, nil
}
