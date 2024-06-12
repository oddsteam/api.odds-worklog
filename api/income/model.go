package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Income struct {
	UserID              string
	NetIncomeStr        string
	NetDailyIncomeStr   string
	NetSpecialIncomeStr string
	VATStr              string
	WHTStr              string
	TotalIncomeStr      string
}

func NewIncome(uidFromSession string) *Income {
	return &Income{
		UserID: uidFromSession,
	}
}

func (i *Income) prepareDataForAddIncome(req models.IncomeReq, userDetail models.User) error {
	ins, err := calIncomeSum(req.WorkDate, userDetail.Vat, userDetail.DailyIncome, userDetail.GetRole())
	if err != nil {
		return err
	}
	insSpecial, err := calIncomeSum(req.WorkingHours, userDetail.Vat, req.SpecialIncome, userDetail.GetRole())
	if err != nil {
		return err
	}
	summaryNetIncome, err := calSummary(ins.Net, insSpecial.Net)
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
	summaryWht, err := calSummary(ins.WHT, insSpecial.WHT)
	if err != nil {
		return err
	}
	summaryIncome, err := calSummary(ins.TotalIncome, insSpecial.TotalIncome)
	if err != nil {
		return err
	}
	i.NetIncomeStr = summaryNetIncome
	i.NetDailyIncomeStr = ins.Net
	i.NetSpecialIncomeStr = insSpecial.Net
	i.VATStr = summaryVat
	i.WHTStr = summaryWht
	i.TotalIncomeStr = summaryIncome
	return nil
}
