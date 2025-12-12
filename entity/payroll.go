package entity

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type Payroll struct {
	UserID            string
	dailyRate         float64
	workDate          float64
	specialHours      float64
	specialIncomeRate float64
	isVATRegistered   bool
	u                 *user.User
	loan              *models.StudentLoan
	data              *models.Income
}

func NewIncome(uidFromSession string) *Payroll {
	return &Payroll{
		UserID: uidFromSession,
		loan:   &models.StudentLoan{},
	}
}

func NewIncomeFromRecord(data models.Income) *Payroll {
	i := Payroll{
		UserID:          "",
		loan:            &models.StudentLoan{},
		data:            &data,
		dailyRate:       data.DailyRate,
		isVATRegistered: data.IsVATRegistered,
	}
	i.parse(IncomeReq{
		SpecialIncome: data.SpecialIncome,
		WorkDate:      data.WorkDate,
		WorkingHours:  data.WorkingHours,
	})
	return &i
}

func CreateIncome(user models.User, req IncomeReq, note string) *models.Income {
	i := NewIncome(user.ID.Hex())
	record, err := i.prepareDataForAddIncome(req, user)
	record.Note = note
	utils.FailOnError(err, "Error prepare data for add income")
	return record
}

func UpdateIncome(user models.User, req IncomeReq, note string, record *models.Income) *models.Income {
	i := NewIncomeFromRecord(*record)
	err := i.prepareDataForUpdateIncome(req, user, record)
	record.Note = note
	utils.FailOnError(err, "Error prepare data for add income")
	return record
}

func (i *Payroll) SetLoan(l *models.StudentLoan) {
	i.loan = l
}

func (i *Payroll) parseRequest(req IncomeReq, userDetail models.User) error {
	err := i.parse(req)
	if err != nil {
		return err
	}
	i.u = user.NewUser(userDetail)
	i.dailyRate = i.u.DailyRate()
	i.isVATRegistered = i.u.IsVATRegistered()
	return nil
}

func (i *Payroll) prepareDataForAddIncome(req IncomeReq, userDetail models.User) (*models.Income, error) {
	income := models.Income{}
	err := i.prepareDataForUpdateIncome(req, userDetail, &income)
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (i *Payroll) prepareDataForUpdateIncome(req IncomeReq, userDetail models.User, income *models.Income) error {
	err := i.parseRequest(req, userDetail)
	if err != nil {
		return err
	}

	income.SubmitDate = time.Now()
	income.UserID = i.UserID
	income.Role = userDetail.Role
	income.ThaiCitizenID = userDetail.ThaiCitizenID
	income.Name = userDetail.GetName()
	income.BankAccountName = userDetail.BankAccountName
	income.BankAccountNumber = userDetail.BankAccountNumber
	income.Email = userDetail.Email
	income.Phone = userDetail.Phone
	income.NetIncome = i.TransferAmountStr()
	income.NetSpecialIncome = i.NetSpecialIncomeStr()
	income.NetDailyIncome = i.NetDailyIncomeStr()
	income.VAT = i.totalVatStr()
	income.WHT = i.TotalWHTStr()
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours
	income.DailyRate = i.u.DailyRate()
	income.IsVATRegistered = i.u.IsVATRegistered()
	income.TotalIncome = i.totalIncomeStr()

	return nil
}

func (i *Payroll) parse(req IncomeReq) error {
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

func (i *Payroll) totalVatStr() string {
	v := i.totalVat()
	if v == 0.0 {
		return ""
	}
	return utils.FloatToString(v)
}

func (i *Payroll) totalVat() float64 {
	return i.VAT(i.totalIncome())
}

func (i *Payroll) totalWHT() float64 {
	return i.WitholdingTax(i.totalIncome())
}

func (i *Payroll) TransferAmountStr() string {
	return utils.FloatToString(i.TransferAmount())
}

func (i *Payroll) TransferAmount() float64 {
	return i.netDailyIncome() + i.netSpecialIncome() - float64(i.loan.Amount)
}

func (i *Payroll) NetDailyIncomeStr() string {
	return utils.FloatToString(i.netDailyIncome())
}

func (i *Payroll) netDailyIncome() float64 {
	return i.Net(i.dailyIncome())
}

func (i *Payroll) totalIncomeStr() string {
	return utils.FloatToString(i.totalIncome())
}

func (i *Payroll) totalIncome() float64 {
	return i.dailyIncome() + i.specialIncome()
}

func (i *Payroll) dailyIncome() float64 {
	return (i.workDate * i.dailyRate)
}

func (i *Payroll) NetSpecialIncomeStr() string {
	return utils.FloatToString(i.netSpecialIncome())
}

func (i *Payroll) netSpecialIncome() float64 {
	return i.Net(i.specialIncome())
}

func (i *Payroll) specialIncome() float64 {
	return i.specialHours * i.specialIncomeRate
}

func (i *Payroll) WitholdingTax(totalIncome float64) float64 {
	return totalIncome * 0.03
}

func (i *Payroll) Net(totalIncome float64) float64 {
	return totalIncome + i.VAT(totalIncome) - i.WitholdingTax(totalIncome)
}

func (i *Payroll) VAT(totalIncome float64) float64 {
	if !i.isVATRegistered {
		return 0
	}
	return totalIncome * 0.07
}

func (i *Payroll) TotalWHTStr() string {
	return utils.FloatToString(i.totalWHT())
}

func (i *Payroll) Note() string {
	return i.data.Note
}

func (i *Payroll) SubmitDateStr() string {
	t := i.data.SubmitDate
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
}

func (i *Payroll) GetName() string {
	return i.data.Name
}

func (i *Payroll) BankAccountNumber() string {
	return i.data.BankAccountNumber
}

func (i *Payroll) ThaiCitizenID() string {
	return i.data.ThaiCitizenID
}

func (i *Payroll) Email() string {
	return i.data.Email
}

func (i *Payroll) GetDeduction() string {
	return i.loan.CSVAmount()
}

func (i *Payroll) GetBankAccountName() string {
	return i.data.BankAccountName
}

func GivenIndividualUser(uidFromSession string, dailyIncome string) models.User {
	return models.User{
		ID:          bson.ObjectIdHex(uidFromSession),
		Role:        "individual",
		Vat:         "N",
		DailyIncome: dailyIncome,
	}
}
