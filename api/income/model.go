package income

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"

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

/** deprecated **/
func (i *Income) export(user models.User) []string {
	income := *i.data
	loan := *i.loan
	t := income.SubmitDate
	netTotalIncome, _ := calTotal(income.NetDailyIncome, income.NetSpecialIncome)
	netTotalIncome = calTotalWithLoanDeduction(netTotalIncome, loan)
	tf := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
	d := []string{
		"",
		user.BankAccountName,
		"",
		utils.SetValueCSV(user.BankAccountNumber),
		user.GetName(),
		user.ThaiCitizenID,
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

const (
	SAP_TXN_INDEX             = 0
	SAP_PAYER_NAME_INDEX      = 1
	SAP_PAYEE_NAME_INDEX      = 2
	SAP_MALE_TO_NAME_INDEX    = 3
	SAP_BENEFICIARY1_INDEX    = 4
	SAP_BENEFICIARY2_INDEX    = 5
	SAP_BENEFICIARY3_INDEX    = 6
	SAP_BENEFICIARY4_INDEX    = 7
	SAP_ZIPCODE_INDEX         = 8
	SAP_CUSTOMER_REF_INDEX    = 9
	SAP_DATE_EFFECTIVE_INDEX  = 10
	SAP_DATE_PICKUP_INDEX     = 11
	SAP_CURRENCY_INDEX        = 12
	SAP_EMPTY_1_INDEX         = 13
	SAP_COMPANY_ACCNO_INDEX   = 14
	SAP_AMOUNT_INDEX          = 15
	SAP_PAYEE_BANK_CODE_INDEX = 16
	SAP_ACCOUNTNO_INDEX       = 17
	SAP_UNKNOW_1_INDEX        = 18
	SAP_EMPTY_2_INDEX         = 19
	SAP_EMPTY_3_INDEX         = 20
	SAP_ADVICEMODE2_INDEX     = 21
	SAP_FAXNO_INDEX           = 22
	SAP_EMAIL_INDEX           = 23
	SAP_SMSNO_INDEX           = 24
	SAP_CHARGE_ON_INDEX       = 25
	SAP_PRODUCT_INDEX         = 26
	SAP_SCHEDULE_INDEX        = 27
	SAP_EMPTY_4_INDEX         = 28
	SAP_DOCREQ_INDEX          = 29
	SAP_EMPTY_5_INDEX         = 30
	SAP_END_INDEX             = 31

	SAP_WHT_WHT_INDEX      = 0
	SAP_WHT_EMPTY_1_INDEX  = 1
	SAP_WHT_TAX_ID_INDEX   = 2
	SAP_WHT_EMPTY_3_INDEX  = 3
	SAP_WHT_EMPTY_4_INDEX  = 4
	SAP_WHT_EMPTY_5_INDEX  = 5
	SAP_WHT_EMPTY_6_INDEX  = 6
	SAP_WHT_EMPTY_7_INDEX  = 7
	SAP_WHT_EMPTY_8_INDEX  = 8
	SAP_WHT_EMPTY_9_INDEX  = 9
	SAP_WHT_EMPTY_10_INDEX = 10
	SAP_WHT_EMPTY_11_INDEX = 11
	SAP_WHT_EMPTY_12_INDEX = 12
	SAP_WHT_EMPTY_13_INDEX = 13
	SAP_WHT_EMPTY_14_INDEX = 14
	SAP_WHT_COM_NAME       = 15
	SAP_WHT_ADDRESS        = 16
	SAP_WHT_EMPTY_17       = 17
	SAP_WHT_EMPTY_18       = 18
	SAP_WHT_EMPTY_19       = 19
	SAP_WHT_EMPTY_20       = 20
)

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

func AddBlank(value string, length int) string {
	if len(value) > length {
		return value[:length]
	}
	return value + strings.Repeat(" ", length-utf8.RuneCountInString(value))
}

func AmountStr(amount float64, n int) string {
	// ปัดทศนิยม 2 หลัก
	rounded := math.Round(amount*100) / 100

	// แปลงเป็น string ด้วย 2 หลักทศนิยม
	s := fmt.Sprintf("%.2f", rounded)

	// ถ้าไม่มีทศนิยม (เช่น 1500 -> 1500.00) Go จะเติมให้แล้ว

	// เติม 0 ด้านหน้าให้ครบความยาว n
	if len(s) > n {
		return s[:n]
	}
	return strings.Repeat("0", n-len(s)) + s
}

func LeftN(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s
}

func ReceiveBRCode(bnkCode, brchCode string) string {
	var tempStr string

	switch bnkCode {
	case "002", "004", "006", "011", "014",
		"015", "020", "022", "024", "025",
		"033", "065", "066", "071", "073":
		// เอา 3 ตัวแรกของสาขาแล้วเติม 0 ข้างหน้า
		if len(brchCode) >= 3 {
			tempStr = "0" + brchCode[:3]
		} else {
			tempStr = "0" + brchCode
		}

	case "030", "067", "072", "069", "052":
		// เอา 4 ตัวแรกของรหัสสาขา
		if len(brchCode) >= 4 {
			tempStr = brchCode[:4]
		} else {
			tempStr = brchCode
		}

	case "034":
		tempStr = "0000"

	case "045":
		tempStr = "0010"

	default:
		tempStr = "0001"
	}

	return tempStr
}

func ReceiveAcCode(bnkCode, acCode string) (string, error) {
	acCode = strings.ReplaceAll(acCode, "-", "")
	acCode = strings.ReplaceAll(acCode, " ", "")
	var tempStr string
	if bnkCode == "011" {
		if len(acCode) == 10 {
			tempStr = "0" + acCode
		} else {
			tempStr = acCode
		}
		return tempStr, nil
	} else {
		return "", errors.New("not supported bank code")
	}

}

func (t Transaction) ToTXNLine() []string {
	bankAccNo, _ := ReceiveAcCode(t.PayYeeBank, t.AccNo)
	return []string{
		"TXN",
		AddBlank(t.ComName, 120),
		AddBlank(t.Payee, 130),
		AddBlank(t.MailTo, 40),
		AddBlank(t.BenAddr1, 40),
		AddBlank(t.BenAddr2, 40),
		AddBlank(t.BenAddr3, 40),
		AddBlank(t.BenAddr4, 40),
		AddBlank(t.ZipCode, 10),
		AddBlank(t.Ref, 16),
		t.DateEff.Format("02012006"), // ddMMyyyy,
		t.DateEff.Format("02012006"), // Date Pick up,
		"THB",
		AddBlank("", 50),
		"0" + AddBlank(t.ComAccNo, 19),
		AmountStr(t.Amt, 15),
		AddBlank(LeftN(t.PayYeeBank, 3), 3) + ReceiveBRCode(t.PayYeeBank, t.AccNo) + AddBlank("", 9),
		AddBlank(bankAccNo, 20),
		"04" + "00",
		AddBlank("", 2),
		AddBlank("", 20), // Pickup Location,
		AddBlank(t.AdvMode, 5),
		AddBlank(strings.ReplaceAll(t.Fax, "-", ""), 50),
		AddBlank(strings.ReplaceAll(t.Mail, "", ""), 50),
		AddBlank(strings.ReplaceAll(t.SMS, "-", ""), 50),
		AddBlank(strings.ToUpper(t.ChargeOn), 13),
		AddBlank(t.Product, 3),
		AddBlank(LeftN(t.Schedule, 5), 5),
		AddBlank("", 34),
		AddBlank(t.DocReq, 105),
		AddBlank("", 295),
		"END",
	}
}

func (t Transaction) ToWHTLine() []string {
	return []string{
		"WHT",
		AddBlank("", 13),
		AddBlank(t.TaxID, 13),
		AddBlank("", 2),
		AmountStr(0, 15),
		AddBlank("", 2),
		AddBlank("", 35),
		AddBlank("", 5),
		AmountStr(0, 15),
		AmountStr(0, 15),
		AddBlank("", 2),
		AddBlank("", 35),
		AddBlank("", 5),
		AmountStr(0, 15),
		AddBlank("", 144),
		AddBlank(t.ComName, 120),
		AddBlank(t.Address, 160),
		AddBlank("", 120),
		AddBlank("", 160),
		AddBlank("", 20),
		AddBlank("", 938),
	}
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

func (i *Income) submitDateStr() string {
	t := i.data.SubmitDate
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
}

type Incomes struct {
	records []*models.Income
	loans   models.StudentLoanList
}

func NewIncomes(records []*models.Income, loans models.StudentLoanList) *Incomes {
	return &Incomes{
		records: records,
		loans:   loans,
	}
}

func NewIncomesWithoutLoans(records []*models.Income) *Incomes {
	return NewIncomes(records, models.StudentLoanList{})
}

func (ics *Incomes) FindByUserID(id string) *models.Income {
	for _, e := range ics.records {
		if id == e.UserID {
			return e
		}
	}
	return &models.Income{}
}

func (ics *Incomes) toCSV() (csv [][]string, updatedIncomeIds []string) {
	strWrite := make([][]string, 0)
	strWrite = append(strWrite, createHeaders())
	updatedIncomeIds = []string{}
	for index, e := range ics.records {
		income := *e
		loan := ics.loans.FindLoan(income.BankAccountName)
		if income.ID.Hex() != "" {
			updatedIncomeIds = append(updatedIncomeIds, income.ID.Hex())
			i := NewIncomeFromRecord(income)
			i.SetLoan(&loan)
			d := i.export2()
			d[VENDOR_CODE_INDEX] = ics.getVendorCode(index)
			strWrite = append(strWrite, d)
		}
	}
	return strWrite, updatedIncomeIds
}

func (ics *Incomes) toCSVasSAP(dateEff time.Time) (csv [][]string, updatedIncomeIds []string) {
	strWrite := make([][]string, 0)
	updatedIncomeIds = []string{}
	for _, e := range ics.records {
		income := *e
		loan := ics.loans.FindLoan(income.BankAccountName)
		if income.ID.Hex() != "" {
			updatedIncomeIds = append(updatedIncomeIds, income.ID.Hex())
			i := NewIncomeFromRecord(income)
			i.SetLoan(&loan)
			txn, wht := i.exportSAP(dateEff)
			strWrite = append(strWrite, txn)
			strWrite = append(strWrite, wht)
		}
	}
	return strWrite, updatedIncomeIds
}

func (ics *Incomes) getVendorCode(i int) string {
	return VendorCode{index: i}.String()
}
func CreateIncome(user models.User, req models.IncomeReq, note string) *models.Income {
	i := NewIncome(string(user.ID))
	record, err := i.prepareDataForAddIncome(req, user)
	record.Note = note
	utils.FailOnError(err, "Error prepare data for add income")
	return record
}

func UpdateIncome(user models.User, req models.IncomeReq, note string, record *models.Income) *models.Income {
	i := NewIncomeFromRecord(*record)
	err := i.prepareDataForUpdateIncome(req, user, record)
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

type VendorCode struct {
	index int
}

func (vc VendorCode) String() string {
	return string([]rune{vc.getFirstLetter(), vc.getSecondLetter(), vc.getThirdLetter()})
}

func (vc VendorCode) getFirstLetter() rune {
	first := 'A' + (vc.index / (26 * 26))
	return rune(first)
}

func (vc VendorCode) getSecondLetter() rune {
	return rune('A' + ((vc.index % (26 * 26)) / 26))
}

func (vc VendorCode) getThirdLetter() rune {
	return rune('A' + (vc.index % 26))
}
