package income

import (
	"fmt"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type Income struct {
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

func NewIncome(uidFromSession string) *Income {
	return &Income{
		UserID: uidFromSession,
		loan:   &models.StudentLoan{},
	}
}

func NewIncomeFromRecord(data models.Income) *Income {
	i := Income{
		UserID:          "",
		loan:            &models.StudentLoan{},
		data:            &data,
		dailyRate:       data.DailyRate,
		isVATRegistered: data.IsVATRegistered,
	}
	i.parse(models.IncomeReq{
		SpecialIncome: data.SpecialIncome,
		WorkDate:      data.WorkDate,
		WorkingHours:  data.WorkingHours,
	})
	return &i
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
	i.dailyRate = i.u.DailyRate
	i.isVATRegistered = i.u.IsVATRegistered()
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
	income.NetIncome = i.transferAmountStr()
	income.NetSpecialIncome = i.netSpecialIncomeStr()
	income.NetDailyIncome = i.netDailyIncomeStr()
	income.VAT = i.summaryVatStr()
	income.WHT = utils.FloatToString(i.summaryWHT())
	income.Note = req.Note
	income.WorkDate = req.WorkDate
	income.SpecialIncome = req.SpecialIncome
	income.WorkingHours = req.WorkingHours
	income.DailyRate = i.u.DailyRate
	income.IsVATRegistered = i.u.IsVATRegistered()
	income.TotalIncome = i.totalIncomeStr()

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

func (i *Income) summaryIncome() float64 {
	return i.totalIncome() + i.specialIncome()
}

func (i *Income) transferAmountStr() string {
	return utils.FloatToString(i.transferAmount())
}

func (i *Income) transferAmount() float64 {
	return i.netDailyIncome() + i.netSpecialIncome() - float64(i.loan.Amount)
}

func (i *Income) netDailyIncomeStr() string {
	return utils.FloatToString(i.netDailyIncome())
}

func (i *Income) netDailyIncome() float64 {
	return i.Net(i.totalIncome())
}

func (i *Income) totalIncomeStr() string {
	return utils.FloatToString(i.totalIncome())
}

func (i *Income) totalIncome() float64 {
	return (i.workDate * i.dailyRate)
}

func (i *Income) netSpecialIncomeStr() string {
	return utils.FloatToString(i.netSpecialIncome())
}

func (i *Income) netSpecialIncome() float64 {
	return i.Net(i.specialIncome())
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
	if !i.isVATRegistered {
		return 0
	}
	return totalIncome * 0.07
}

func (i *Income) export(user models.User) []string {
	income := *i.data
	loan := *i.loan
	t := income.SubmitDate
	summaryIncome, _ := calSummary(income.NetDailyIncome, income.NetSpecialIncome)
	summaryIncome = calSummaryWithLoan(summaryIncome, loan)
	tf := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
	d := []string{
		user.GetName(),
		user.ThaiCitizenID,
		user.BankAccountName,
		utils.SetValueCSV(user.BankAccountNumber),
		user.Email,
		utils.FormatCommas(income.NetDailyIncome),
		utils.FormatCommas(income.NetSpecialIncome),
		loan.CSVAmount(),
		income.WHT,
		utils.FormatCommas(summaryIncome),
		income.Note,
		tf,
	}
	return d
}

func calSummaryWithLoan(summaryIncome string, loan models.StudentLoan) string {
	summary, _ := utils.StringToFloat64(summaryIncome)
	summary = summary - float64(loan.Amount)
	summaryIncome = utils.FloatToString(summary)
	return summaryIncome
}

func (i *Income) export2() []string {
	income := *i.data
	loan := *i.loan
	d := []string{
		income.Name,
		income.ThaiCitizenID,
		income.BankAccountName,
		utils.SetValueCSV(income.BankAccountNumber),
		income.Email,
		utils.FormatCommas(income.NetDailyIncome),
		utils.FormatCommas(income.NetSpecialIncome),
		loan.CSVAmount(),
		income.WHT,
		utils.FormatCommas(i.transferAmountStr()),
		income.Note,
		i.submitDateStr(),
	}
	return d
}

func (i *Income) submitDateStr() string {
	t := i.data.SubmitDate
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
}

type Incomes struct {
	users   []*models.User
	records []*models.Income
	loans   models.StudentLoanList
}

func NewIncomes(records []*models.Income, loans models.StudentLoanList, users []*models.User) *Incomes {
	return &Incomes{
		records: records,
		loans:   loans,
		users:   users,
	}
}

func (ics *Incomes) toCSV() (csv [][]string, updatedIncomeIds []string) {
	strWrite := make([][]string, 0)
	strWrite = append(strWrite, createHeaders())
	updatedIncomeIds = []string{}
	for _, user := range ics.users {
		income := models.Income{}
		for _, e := range ics.records {
			if strings.Contains(user.ID.Hex(), e.UserID) {
				income = *e
			}
		}
		loan := ics.loans.FindLoan(*user)
		if income.ID.Hex() != "" {
			updatedIncomeIds = append(updatedIncomeIds, income.ID.Hex())
			i := NewIncomeFromRecord(income)
			i.SetLoan(&loan)
			d := i.export2()
			strWrite = append(strWrite, d)
		}
	}
	return strWrite, updatedIncomeIds
}
