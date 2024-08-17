package income

import (
	"fmt"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
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
	income.Role = userDetail.Role
	income.ThaiCitizenID = userDetail.ThaiCitizenID
	income.Name = userDetail.GetName()
	income.BankAccountName = userDetail.BankAccountName
	income.BankAccountNumber = userDetail.BankAccountNumber
	income.Email = userDetail.Email
	income.Phone = userDetail.Phone
	income.NetIncome = i.transferAmountStr()
	income.NetSpecialIncome = i.netSpecialIncomeStr()
	income.NetDailyIncome = i.netDailyIncomeStr()
	income.VAT = i.totalVatStr()
	income.WHT = utils.FloatToString(i.totalWHT())
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

func (i *Income) totalVatStr() string {
	v := i.totalVat()
	if v == 0.0 {
		return ""
	}
	return utils.FloatToString(v)
}

func (i *Income) totalVat() float64 {
	return i.VAT(i.totalIncome())
}

func (i *Income) totalWHT() float64 {
	return i.WitholdingTax(i.totalIncome())
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
	return i.Net(i.dailyIncome())
}

func (i *Income) totalIncomeStr() string {
	return utils.FloatToString(i.totalIncome())
}

func (i *Income) totalIncome() float64 {
	return i.dailyIncome() + i.specialIncome()
}

func (i *Income) dailyIncome() float64 {
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
	netTotalIncome, _ := calTotal(income.NetDailyIncome, income.NetSpecialIncome)
	netTotalIncome = calTotalWithLoanDeduction(netTotalIncome, loan)
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
		utils.FormatCommas(netTotalIncome),
		income.Note,
		tf,
	}
	return d
}

func calTotalWithLoanDeduction(totalIncomeStr string, loan models.StudentLoan) string {
	totalIncome, _ := utils.StringToFloat64(totalIncomeStr)
	totalIncome = totalIncome - float64(loan.Amount)
	totalIncomeStr = utils.FloatToString(totalIncome)
	return totalIncomeStr
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
		loan := ics.loans.FindLoan(income.BankAccountName)
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

func CreateIncome(user models.User, req models.IncomeReq, note string) *models.Income {
	i := NewIncome(string(user.ID))
	record, err := i.prepareDataForAddIncome(req, user)
	record.Note = note
	utils.FailOnError(err, "Error prepare data for add income")
	return record
}

func GivenIndividualUser(uidFromSession string, dailyIncome string) models.User {
	return models.User{
		ID:          bson.ObjectIdHex(uidFromSession),
		Role:        "individual",
		Vat:         "N",
		DailyIncome: dailyIncome,
	}
}
