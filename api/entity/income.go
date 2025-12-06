package entity

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

func (i *Income) SetLoan(l *models.StudentLoan) {
	i.loan = l
}

func (i *Income) parseRequest(req IncomeReq, userDetail models.User) error {
	err := i.parse(req)
	if err != nil {
		return err
	}
	i.u = user.NewUser(userDetail)
	i.dailyRate = i.u.DailyRate()
	i.isVATRegistered = i.u.IsVATRegistered()
	return nil
}

func (i *Income) prepareDataForAddIncome(req IncomeReq, userDetail models.User) (*models.Income, error) {
	income := models.Income{}
	err := i.prepareDataForUpdateIncome(req, userDetail, &income)
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (i *Income) prepareDataForUpdateIncome(req IncomeReq, userDetail models.User, income *models.Income) error {
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
	income.DailyRate = i.u.DailyRate()
	income.IsVATRegistered = i.u.IsVATRegistered()
	income.TotalIncome = i.totalIncomeStr()

	return nil
}

func (i *Income) parse(req IncomeReq) error {
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

func (i *Income) export() []string {
	income := *i.data
	loan := *i.loan
	d := []string{
		"",
		income.BankAccountName,
		"",
		utils.SetValueCSV(income.BankAccountNumber),
		income.Name,
		income.ThaiCitizenID,
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

func (t *Income) toTransaction(dateEff time.Time) Transaction {
	return Transaction{
		ComName:    "บจก. ออด-อี (ประเทศไทย) จำกัด",
		Payee:      t.data.Name,
		MailTo:     "",
		BenAddr1:   "",
		BenAddr2:   "",
		BenAddr3:   "",
		BenAddr4:   "",
		ZipCode:    "",
		Ref:        "",
		DateEff:    dateEff,
		ComAccNo:   "0011595873",
		Amt:        t.transferAmount(),
		PayYeeBank: "011",
		RecBRCode:  "",
		AccNo:      t.data.BankAccountNumber,
		AdvMode:    "",
		Fax:        "",
		Mail:       "",
		SMS:        "",
		ChargeOn:   "OUR",
		Product:    "DCR",
		Schedule:   "",
		DocReq:     "",
		Address:    "2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900",
		TaxID:      "0105556110718",
	}
}
func (i *Income) exportSAP(dateEff time.Time) ([]string, []string) {
	txn := i.toTransaction(dateEff)
	return txn.ToTXNLine(), txn.ToWHTLine()
}

type Transaction struct {
	ComName    string
	Payee      string
	MailTo     string
	BenAddr1   string
	BenAddr2   string
	BenAddr3   string
	BenAddr4   string
	ZipCode    string
	Ref        string
	DateEff    time.Time
	ComAccNo   string
	Amt        float64
	PayYeeBank string
	RecBRCode  string
	AccNo      string
	AdvMode    string
	Fax        string
	Mail       string
	SMS        string
	ChargeOn   string
	Product    string
	Schedule   string
	DocReq     string
	Address    string
	TaxID      string
}

func (t Transaction) ToTXNLine() []string {
	bankAccNo, _ := utils.ReceiveAcCode(t.PayYeeBank, t.AccNo)
	return []string{
		"TXN",
		utils.AddBlank(t.ComName, 120),
		utils.AddBlank(utils.FilterOthersThanThaiAndAscii(t.Payee), 130),
		utils.AddBlank(t.MailTo, 40),
		utils.AddBlank(t.BenAddr1, 40),
		utils.AddBlank(t.BenAddr2, 40),
		utils.AddBlank(t.BenAddr3, 40),
		utils.AddBlank(t.BenAddr4, 40),
		utils.AddBlank(t.ZipCode, 10),
		utils.AddBlank(t.Ref, 16),
		t.DateEff.Format("02012006"), // ddMMyyyy,
		t.DateEff.Format("02012006"), // Date Pick up,
		"THB",
		utils.AddBlank("", 50),
		"0" + utils.AddBlank(t.ComAccNo, 19),
		utils.AmountStr(t.Amt, 15),
		utils.AddBlank(utils.LeftN(t.PayYeeBank, 3), 3) + utils.ReceiveBRCode(t.PayYeeBank, t.AccNo) + utils.AddBlank("", 9),
		utils.AddBlank(bankAccNo, 20),
		"04" + "00",
		utils.AddBlank("", 2),
		utils.AddBlank("", 20), // Pickup Location,
		utils.AddBlank(t.AdvMode, 5),
		utils.AddBlank(strings.ReplaceAll(t.Fax, "-", ""), 50),
		utils.AddBlank(strings.ReplaceAll(t.Mail, "", ""), 50),
		utils.AddBlank(strings.ReplaceAll(t.SMS, "-", ""), 50),
		utils.AddBlank(strings.ToUpper(t.ChargeOn), 13),
		utils.AddBlank(t.Product, 3),
		utils.AddBlank(utils.LeftN(t.Schedule, 5), 5),
		utils.AddBlank("", 34),
		utils.AddBlank(t.DocReq, 105),
		utils.AddBlank("", 295),
		"END",
	}
}

func (t Transaction) ToWHTLine() []string {
	return []string{
		"WHT",
		utils.AddBlank("", 13),
		utils.AddBlank(t.TaxID, 13),
		utils.AddBlank("", 2),
		utils.AmountStr(0, 15),
		utils.AddBlank("", 2),
		utils.AddBlank("", 35),
		utils.AddBlank("", 5),
		utils.AmountStr(0, 15),
		utils.AmountStr(0, 15),
		utils.AddBlank("", 2),
		utils.AddBlank("", 35),
		utils.AddBlank("", 5),
		utils.AmountStr(0, 15),
		utils.AddBlank("", 144),
		utils.AddBlank(t.ComName, 120),
		utils.AddBlank(t.Address, 160),
		utils.AddBlank("", 120),
		utils.AddBlank("", 160),
		utils.AddBlank("", 20),
		utils.AddBlank("", 938),
	}
}

func GivenIndividualUser(uidFromSession string, dailyIncome string) models.User {
	return models.User{
		ID:          bson.ObjectIdHex(uidFromSession),
		Role:        "individual",
		Vat:         "N",
		DailyIncome: dailyIncome,
	}
}
