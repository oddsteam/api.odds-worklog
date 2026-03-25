package file

import (
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
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

func toTransaction(i models.Payroll, dateEff time.Time) Transaction {
	return Transaction{
		ComName:    "บจก. ออด-อี (ประเทศไทย) จำกัด",
		Payee:      i.GetName(),
		MailTo:     "",
		BenAddr1:   "",
		BenAddr2:   "",
		BenAddr3:   "",
		BenAddr4:   "",
		ZipCode:    "",
		Ref:        "",
		DateEff:    dateEff,
		ComAccNo:   "0011595873",
		Amt:        i.TransferAmount(),
		PayYeeBank: "011",
		RecBRCode:  "",
		AccNo:      i.BankAccountNumber(),
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
