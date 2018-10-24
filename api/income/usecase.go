package income

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

func newUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
}

func calVAT(income string) (string, error) {
	num, err := strconv.ParseFloat(income, 64)
	if err != nil {
		return "", err
	}
	vat := float64(num * 0.07)
	return fmt.Sprintf("%f", vat), nil
}

func calWHT(income string) (string, error) {
	num, err := strconv.ParseFloat(income, 64)
	if err != nil {
		return "", err
	}
	vat := float64(num * 0.03)
	return fmt.Sprintf("%f", vat), nil
}

func (u *usecase) AddIncome(req *models.IncomeReq, id string) error {
	println(id)
	us, err := u.userRepo.GetUserByID(id)
	if err != nil {
		println(err)
		return err
	}

	var vat, wht = "", ""
	vat, err = calVAT(req.TotalIncome)
	if err != nil {
		return err
	}
	if us.CorporateFlag == "Y" {
		wht, err = calWHT(req.TotalIncome)
		if err != nil {
			return err
		}
	}

	t := time.Now()

	income := models.Income{
		UserID:      id,
		TotalIncome: req.TotalIncome,
		SubmitDate:  t.Format(time.RFC3339),
		Note:        req.Note,
		VAT:         vat,
		WHT:         wht,
	}
	err = u.repo.AddIncome(&income)
	if err != nil {
		return err
	}
	return nil
}
