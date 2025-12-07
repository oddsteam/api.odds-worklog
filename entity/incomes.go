package entity

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

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

func (ics *Incomes) processRecords(process func(index int, i *Income) [][]string) ([][]string, []string) {
	strWrite := make([][]string, 0)
	updatedIncomeIds := []string{}
	for index, e := range ics.records {
		income := *e
		if income.ID.Hex() != "" {
			updatedIncomeIds = append(updatedIncomeIds, income.ID.Hex())
			loan := ics.loans.FindLoan(income.BankAccountName)
			i := NewIncomeFromRecord(income)
			i.SetLoan(&loan)
			rows := process(index, i)
			strWrite = append(strWrite, rows...)
		}
	}
	return strWrite, updatedIncomeIds
}

func (ics *Incomes) getVendorCode(i int) string {
	return VendorCode{index: i}.String()
}

func (ics *Incomes) ToCSV() ([][]string, []string) {
	rows, ids := ics.processRecords(func(index int, i *Income) [][]string {
		d := i.export()
		d[VENDOR_CODE_INDEX] = ics.getVendorCode(index)
		return [][]string{d}
	})
	return append([][]string{createHeaders()}, rows...), ids
}

func (ics *Incomes) ToSAP(dateEff time.Time) ([][]string, []string) {
	return ics.processRecords(func(index int, i *Income) [][]string {
		txn, wht := i.exportSAP(dateEff)
		return [][]string{txn, wht}
	})
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

const (
	VENDOR_CODE_INDEX = iota
	ACCOUNT_NAME_INDEX
	PAYMENT_METHOD_INDEX
	ACCOUNT_NUMBER_INDEX
	NAME_INDEX
	ID_CARD_INDEX
	EMAIL_INDEX
	NET_DAILY_INCOME_INDEX
	NET_SPECIAL_INCOME_INDEX
	LOAN_DEDUCTION_INDEX
	WITHHOLDING_TAX_INDEX
	TRANSFER_AMOUNT_INDEX
	NOTE_INDEX
	SUBMIT_DATE_INDEX
)

func createHeaders() []string {
	return []string{"Vendor Code", "ชื่อบัญชี", "Payment method", "เลขบัญชี", "ชื่อ", "เลขบัตรประชาชน", "อีเมล", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ", "หัก ณ ที่จ่าย", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
}
