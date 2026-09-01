package models

import (
	"fmt"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

const (
	CurrentWHTRate = 0.05
	LegacyWHTRate  = 0.03
)

type Payroll struct {
	UserID              string
	dailyRate           float64
	workedDays          float64
	specialWorkingHours float64
	specialHourlyRate   float64
	isVATRegistered     bool
	whtRate             float64
	userProfile         *User
	loan                *StudentLoan
	incomeRecord        *Income
}

func NewPayroll(uidFromSession string) *Payroll {
	return &Payroll{
		UserID:  uidFromSession,
		loan:    &StudentLoan{},
		whtRate: CurrentWHTRate,
	}
}

func NewPayrollFromIncome(record Income) *Payroll {
	p := Payroll{
		UserID:          record.UserID,
		loan:            &StudentLoan{},
		incomeRecord:    &record,
		dailyRate:       record.DailyRate,
		isVATRegistered: record.IsVATRegistered,
		whtRate:         record.EffectiveWHTRate(),
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
	i.whtRate = CurrentWHTRate
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
	p.userProfile = &userDetail
	p.dailyRate = p.getUserDailyRate()
	p.isVATRegistered = p.isUserVATRegistered()
	return nil
}

func (p *Payroll) getUserDailyRate() float64 {
	dr, _ := StringToFloat64(p.userProfile.DailyIncome)
	return dr
}

func (p *Payroll) isUserVATRegistered() bool {
	return p.userProfile.Vat == "Y"
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
	income.PeakCode = userDetail.PeakCode
	income.NetIncome = p.TransferAmountStr()
	income.NetSpecialIncome = p.NetSpecialIncomeStr()
	income.DailyIncomeBeforeTax = p.DailyIncomeBeforeTaxStr()
	income.SpecialIncomeBeforeTax = p.SpecialIncomeBeforeTaxStr()
	income.NetDailyIncome = p.NetDailyIncomeStr()
	income.VAT = p.totalVatStr()
	income.WHT = p.TotalWHTStr()
	income.WHTRate = p.whtRate
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
	p.workedDays, err = StringToFloat64(req.WorkDate)
	if err != nil {
		p.workedDays = 0
	}
	p.specialWorkingHours, err = StringToFloat64(req.WorkingHours)
	if err != nil {
		p.specialWorkingHours = 0
	}
	p.specialHourlyRate, err = StringToFloat64(req.SpecialIncome)
	if err != nil {
		p.specialHourlyRate = 0
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

func (p *Payroll) DailyIncomeBeforeTaxStr() string {
	return FloatToString(p.dailyIncome())
}

func (p *Payroll) SpecialIncomeBeforeTaxStr() string {
	return FloatToString(p.specialIncome())
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
	return (p.workedDays * p.dailyRate)
}

func (p *Payroll) NetSpecialIncomeStr() string {
	return FloatToString(p.netSpecialIncome())
}

func (p *Payroll) netSpecialIncome() float64 {
	return p.Net(p.specialIncome())
}

func (p *Payroll) specialIncome() float64 {
	return p.specialWorkingHours * p.specialHourlyRate
}

func (p *Payroll) WitholdingTax(totalIncome float64) float64 {
	return totalIncome * p.whtRate
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

func (p *Payroll) TotalVATStr() string {
	return p.totalVatStr()
}

func (p *Payroll) Note() string {
	return p.incomeRecord.Note
}

func (p *Payroll) SubmitDateStr() string {
	t := p.incomeRecord.SubmitDate
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
}

func (p *Payroll) GetName() string {
	return p.incomeRecord.Name
}

func (p *Payroll) SiteName() string {
	if p.incomeRecord == nil {
		return ""
	}
	return p.incomeRecord.SiteName
}

func (p *Payroll) BankAccountNumber() string {
	return p.incomeRecord.BankAccountNumber
}

func (p *Payroll) ThaiCitizenID() string {
	return p.incomeRecord.ThaiCitizenID
}

func (p *Payroll) Email() string {
	return p.incomeRecord.Email
}

func (p *Payroll) WorkDate() string {
	return p.incomeRecord.WorkDate
}

func (p *Payroll) WorkingHours() string {
	return p.incomeRecord.WorkingHours
}

func (p *Payroll) SpecialIncomeRateStr() string {
	return p.incomeRecord.SpecialIncome
}

func (p *Payroll) GetDeduction() string {
	return p.loan.CSVAmount()
}

func (p *Payroll) GetBankAccountName() string {
	return p.incomeRecord.BankAccountName
}

func GivenIndividualUser(uidFromSession string, dailyIncome string) User {
	return User{
		ID:          bsonutil.MustObjectIDFromHex(uidFromSession),
		Role:        "individual",
		Vat:         "N",
		DailyIncome: dailyIncome,
	}
}
