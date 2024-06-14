package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type Income struct {
	UserID            string
	workDate          float64
	specialHours      float64
	specialIncomeRate float64
	u                 *user.User
	loan              *models.StudentLoan
}

func NewIncome(uidFromSession string) *Income {
	return &Income{
		UserID: uidFromSession,
		loan:   &models.StudentLoan{},
	}
}

func (i *Income) SetLoan(l *models.StudentLoan) {
	i.loan = l
}

func (i *Income) parseRequest(req models.IncomeReq, userDetail models.User) error {
	err := i.parse(req)
	if err != nil {
		return err
	}
	i.u = user.NewUser(userDetail)
	i.u.Parse()
	return nil
}

func (i *Income) prepareDataForAddIncome(req models.IncomeReq, userDetail models.User) (*models.Income, error) {
	income := models.Income{}
	err := i.prepareDataForUpdateIncome(req, userDetail, &income)
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (i *Income) prepareDataForUpdateIncome(req models.IncomeReq, userDetail models.User, income *models.Income) error {
	err := i.parseRequest(req, userDetail)
	if err != nil {
		return err
	}

	income.SubmitDate = time.Now()
	income.UserID = i.UserID
	income.ThaiCitizenID = userDetail.ThaiCitizenID
	income.Name = userDetail.GetName()
	income.BankAccountName = userDetail.BankAccountName
	income.BankAccountNumber = userDetail.BankAccountNumber
	income.Email = userDetail.Email
	income.Phone = userDetail.Phone
	income.TotalIncome = utils.FloatToString(i.summaryIncome())
	income.NetIncome = utils.FloatToString(i.summaryNet())
	income.NetSpecialIncome = utils.FloatToString(i.Net(i.specialIncome()))
	income.NetDailyIncome = utils.FloatToString(i.Net(i.totalIncome()))
	income.VAT = i.summaryVatStr()
	income.WHT = utils.FloatToString(i.summaryWHT())
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours

	return nil
}

func (i *Income) parse(req models.IncomeReq) error {
	var err error
	i.workDate, err = utils.StringToFloat64(req.WorkDate)
	if err != nil {
		i.workDate = 0
	}
	i.specialHours, err = utils.StringToFloat64(req.WorkingHours)
	if err != nil {
		i.specialHours = 0
	}
	i.specialIncomeRate, err = utils.StringToFloat64(req.SpecialIncome)
	if err != nil {
		i.specialIncomeRate = 0
	}
	return nil
}

func (i *Income) summaryVatStr() string {
	v := i.summaryVat()
	if v == 0.0 {
		return ""
	}
	return utils.FloatToString(v)
}
func (i *Income) summaryVat() float64 {
	return i.VAT(i.summaryIncome())
}

func (i *Income) summaryWHT() float64 {
	return i.WitholdingTax(i.summaryIncome())
}

func (i *Income) summaryNet() float64 {
	return i.Net(i.summaryIncome())
}

func (i *Income) summaryIncome() float64 {
	return i.totalIncome() + i.specialIncome()
}

func (i *Income) totalIncome() float64 {
	return (i.workDate * i.u.DailyRate) - float64(i.loan.Amount)
}

func (i *Income) specialIncome() float64 {
	return i.specialHours * i.specialIncomeRate
}

func (i *Income) WitholdingTax(totalIncome float64) float64 {
	return totalIncome * 0.03
}

func (i *Income) Net(totalIncome float64) float64 {
	return totalIncome + i.VAT(totalIncome) - i.WitholdingTax(totalIncome)
}
func (i *Income) VAT(totalIncome float64) float64 {
	if !i.u.IsVATRegistered() {
		return 0
	}
	return totalIncome * 0.07
}
