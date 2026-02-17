package models

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

type Payroll struct {
	UserID            string
	dailyRate         float64
	workDate          float64
	specialHours      float64
	specialIncomeRate float64
	isVATRegistered   bool
	userDetail        *User
	loan              *StudentLoan
	record            *Income
}

func NewPayroll(uidFromSession string) *Payroll {
	return &Payroll{
		UserID: uidFromSession,
		loan:   &StudentLoan{},
	}
}

func NewPayrollFromIncome(record Income) *Payroll {
	p := Payroll{
		UserID:          "",
		loan:            &StudentLoan{},
		record:          &record,
		dailyRate:       record.DailyRate,
		isVATRegistered: record.IsVATRegistered,
	}
	p.parse(IncomeReq{
		SpecialIncome: record.SpecialIncome,
		WorkDate:      record.WorkDate,
		WorkingHours:  record.WorkingHours,
	})
	return &p
}

func CreatePayroll(user User, req IncomeReq, note string) *Income {
	i := NewPayroll(user.ID.Hex())
	record, err := i.prepareDataForAddIncome(req, user)
	record.Note = note
	FailOnError(err, "Error prepare data for add income")
	return record
}

func UpdatePayroll(user User, req IncomeReq, note string, record *Income) *Income {
	i := NewPayrollFromIncome(*record)
	err := i.prepareDataForUpdateIncome(req, user, record)
	record.Note = note
	FailOnError(err, "Error prepare data for add income")
	return record
}

func (p *Payroll) SetLoan(l *StudentLoan) {
	p.loan = l
}

func (p *Payroll) parseRequest(req IncomeReq, userDetail User) error {
	err := p.parse(req)
	if err != nil {
		return err
	}
	p.userDetail = &userDetail
	p.dailyRate = p.getUserDailyRate()
	p.isVATRegistered = p.isUserVATRegistered()
	return nil
}

func (p *Payroll) getUserDailyRate() float64 {
	dr, _ := StringToFloat64(p.userDetail.DailyIncome)
	return dr
}

func (p *Payroll) isUserVATRegistered() bool {
	return p.userDetail.Vat == "Y"
}

func (p *Payroll) prepareDataForAddIncome(req IncomeReq, userDetail User) (*Income, error) {
	income := Income{}
	err := p.prepareDataForUpdateIncome(req, userDetail, &income)
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (p *Payroll) prepareDataForUpdateIncome(req IncomeReq, userDetail User, income *Income) error {
	err := p.parseRequest(req, userDetail)
	if err != nil {
		return err
	}

	income.SubmitDate = time.Now()
	income.UserID = p.UserID
	income.Role = userDetail.Role
	income.ThaiCitizenID = userDetail.ThaiCitizenID
	income.Name = userDetail.GetName()
	income.BankAccountName = userDetail.BankAccountName
	income.BankAccountNumber = userDetail.BankAccountNumber
	income.Email = userDetail.Email
	income.Phone = userDetail.Phone
	income.NetIncome = p.TransferAmountStr()
	income.NetSpecialIncome = p.NetSpecialIncomeStr()
	income.NetDailyIncome = p.NetDailyIncomeStr()
	income.VAT = p.totalVatStr()
	income.WHT = p.TotalWHTStr()
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours
	income.DailyRate = p.getUserDailyRate()
	income.IsVATRegistered = p.isUserVATRegistered()
	income.TotalIncome = p.totalIncomeStr()

	return nil
}

func (p *Payroll) parse(req IncomeReq) error {
	var err error
	p.workDate, err = StringToFloat64(req.WorkDate)
	if err != nil {
		p.workDate = 0
	}
	p.specialHours, err = StringToFloat64(req.WorkingHours)
	if err != nil {
		p.specialHours = 0
	}
	p.specialIncomeRate, err = StringToFloat64(req.SpecialIncome)
	if err != nil {
		p.specialIncomeRate = 0
	}
	return nil
}

func (p *Payroll) totalVatStr() string {
	v := p.totalVat()
	if v == 0.0 {
		return ""
	}
	return FloatToString(v)
}

func (p *Payroll) totalVat() float64 {
	return p.VAT(p.totalIncome())
}

func (p *Payroll) totalWHT() float64 {
	return p.WitholdingTax(p.totalIncome())
}

func (p *Payroll) TransferAmountStr() string {
	return FloatToString(p.TransferAmount())
}

func (p *Payroll) TransferAmount() float64 {
	return p.netDailyIncome() + p.netSpecialIncome() - float64(p.loan.Amount)
}

func (p *Payroll) NetDailyIncomeStr() string {
	return FloatToString(p.netDailyIncome())
}

func (p *Payroll) netDailyIncome() float64 {
	return p.Net(p.dailyIncome())
}

func (p *Payroll) totalIncomeStr() string {
	return FloatToString(p.totalIncome())
}

func (p *Payroll) totalIncome() float64 {
	return p.dailyIncome() + p.specialIncome()
}

func (p *Payroll) dailyIncome() float64 {
	return (p.workDate * p.dailyRate)
}

func (p *Payroll) NetSpecialIncomeStr() string {
	return FloatToString(p.netSpecialIncome())
}

func (p *Payroll) netSpecialIncome() float64 {
	return p.Net(p.specialIncome())
}

func (p *Payroll) specialIncome() float64 {
	return p.specialHours * p.specialIncomeRate
}

func (p *Payroll) WitholdingTax(totalIncome float64) float64 {
	return totalIncome * 0.03
}

func (p *Payroll) Net(totalIncome float64) float64 {
	return totalIncome + p.VAT(totalIncome) - p.WitholdingTax(totalIncome)
}

func (p *Payroll) VAT(totalIncome float64) float64 {
	if !p.isVATRegistered {
		return 0
	}
	return totalIncome * 0.07
}

func (p *Payroll) TotalWHTStr() string {
	return FloatToString(p.totalWHT())
}

func (p *Payroll) Note() string {
	return p.record.Note
}

func (p *Payroll) SubmitDateStr() string {
	t := p.record.SubmitDate
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
}

func (p *Payroll) GetName() string {
	return p.record.Name
}

func (p *Payroll) BankAccountNumber() string {
	return p.record.BankAccountNumber
}

func (p *Payroll) ThaiCitizenID() string {
	return p.record.ThaiCitizenID
}

func (p *Payroll) Email() string {
	return p.record.Email
}

func (p *Payroll) GetDeduction() string {
	return p.loan.CSVAmount()
}

func (p *Payroll) GetBankAccountName() string {
	return p.record.BankAccountName
}

func GivenIndividualUser(uidFromSession string, dailyIncome string) User {
	return User{
		ID:          bson.ObjectIdHex(uidFromSession),
		Role:        "individual",
		Vat:         "N",
		DailyIncome: dailyIncome,
	}
}
