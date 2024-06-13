package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Income struct {
	UserID string
}

func NewIncome(uidFromSession string) *Income {
	return &Income{
		UserID: uidFromSession,
	}
}

func (i *Income) prepareDataForAddIncome(req models.IncomeReq, userDetail models.User) (*models.Income, error) {
	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome, userDetail.GetRole())
	if err != nil {
		return nil, err
	}
	insSpecial, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome, userDetail.GetRole())
	if err != nil {
		return nil, err
	}
	summaryNetIncome, err := calSummary(ins.Net, insSpecial.Net)
	if err != nil {
		return nil, err
	}
	var summaryVat string
	if userDetail.Vat != "N" {
		summaryVat, err = calSummary(ins.VAT, insSpecial.VAT)
		if err != nil {
			return nil, err
		}
	} else {
		summaryVat = ""
	}
	summaryWht, err := calSummary(ins.WHT, insSpecial.WHT)
	if err != nil {
		return nil, err
	}
	summaryIncome, err := calSummary(ins.TotalIncome, insSpecial.TotalIncome)
	if err != nil {
		return nil, err
	}

	income := models.Income{
		UserID:            i.UserID,
		TotalIncome:       summaryIncome,
		NetIncome:         summaryNetIncome,
		NetSpecialIncome:  insSpecial.Net,
		NetDailyIncome:    ins.Net,
		Note:              req.Note,
		VAT:               summaryVat,
		WHT:               summaryWht,
		WorkDate:          req.WorkDate,
		SpecialIncome:     req.SpecialIncome,
		WorkingHours:      req.WorkingHours,
		ThaiCitizenID:     userDetail.ThaiCitizenID,
		Name:              userDetail.GetName(),
		BankAccountName:   userDetail.BankAccountName,
		BankAccountNumber: userDetail.BankAccountNumber,
		Email:             userDetail.Email,
		Phone:             userDetail.Phone,
	}

	return &income, nil
}

func (i *Income) prepareDataForUpdateIncome(req models.IncomeReq, userDetail models.User, income *models.Income) error {
	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome, userDetail.GetRole())
	if err != nil {
		return err
	}
	insSpecial, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome, userDetail.GetRole())
	if err != nil {
		return err
	}
	summaryIncome, err := calSummary(ins.TotalIncome, insSpecial.TotalIncome)
	if err != nil {
		return err
	}
	summaryNetIncome, err := calSummary(ins.Net, insSpecial.Net)
	if err != nil {
		return err
	}
	summaryWht, err := calSummary(ins.WHT, insSpecial.WHT)
	if err != nil {
		return err
	}
	var summaryVat string
	if userDetail.Vat != "N" {
		summaryVat, err = calSummary(ins.VAT, insSpecial.VAT)
		if err != nil {
			return err
		}
	} else {
		summaryVat = ""
	}

	income.SubmitDate = time.Now()
	income.TotalIncome = summaryIncome
	income.NetIncome = summaryNetIncome
	income.NetSpecialIncome = insSpecial.Net
	income.NetDailyIncome = ins.Net
	income.VAT = summaryVat
	income.WHT = summaryWht
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours

	return nil
}
