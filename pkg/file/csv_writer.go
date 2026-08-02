package file

import (
	"encoding/csv"
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type csvWriter struct{}

func NewCSVWriter() *csvWriter {
	return &csvWriter{}
}

func (w *csvWriter) WriteFile(name string, ics models.PayrollCycle) (string, error) {
	strWrite, _ := ToCSV(ics)

	if len(strWrite) == 1 {
		return "", errors.New("no data for export to CSV file")
	}

	file, filename, err := CreateFile(name)

	if err != nil {
		return "", err
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()
	defer file.Close()
	return filename, nil
}

func ToCSV(ics models.PayrollCycle) ([][]string, []string) {
	rows, metas := ics.ProcessRecords(func(index int, i models.Payroll) [][]string {
		d := export(i)
		d[VENDOR_CODE_INDEX] = getVendorCode(index)
		return [][]string{d}
	})
	ids := make([]string, len(metas))
	for i, m := range metas {
		ids[i] = m.IncomeID
	}
	return append([][]string{createHeaders()}, rows...), ids
}

func getVendorCode(i int) string {
	return VendorCode{index: i}.String()
}

const (
	VENDOR_CODE_INDEX = iota
	ACCOUNT_NAME_INDEX
	PAYMENT_METHOD_INDEX
	ACCOUNT_NUMBER_INDEX
	NAME_INDEX
	SITE_INDEX
	ID_CARD_INDEX
	EMAIL_INDEX
	WORK_DATE_INDEX
	DAILY_INCOME_BEFORE_TAX_INDEX
	NET_DAILY_INCOME_INDEX
	WORKING_HOURS_INDEX
	HOURLY_RATE_INDEX
	SPECIAL_INCOME_BEFORE_TAX_INDEX
	NET_SPECIAL_INCOME_INDEX
	LOAN_DEDUCTION_INDEX
	VAT_INDEX
	WITHHOLDING_TAX_INDEX
	TRANSFER_AMOUNT_INDEX
	NOTE_INDEX
	SUBMIT_DATE_INDEX
)

func createHeaders() []string {
	return []string{"Vendor Code", "ชื่อบัญชี", "Payment method", "เลขบัญชี", "ชื่อ", "Site", "เลขบัตรประชาชน", "อีเมล", "จำนวนวันทำงาน", "รายได้หลักก่อนคำนวนภาษี", "จำนวนเงินรายได้หลัก", "จำนวนชั่วโมงงานพิเศษ", "อัตราต่อชั่วโมง", "รายได้พิเศษก่อนคำนวนภาษี", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ", "VAT", "หัก ณ ที่จ่าย", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
}

func export(i models.Payroll) []string {
	d := []string{
		"",
		i.GetBankAccountName(),
		"",
		utils.SetValueCSV(i.BankAccountNumber()),
		i.GetName(),
		i.SiteName(),
		i.ThaiCitizenID(),
		i.Email(),
		i.WorkDate(),
		models.FormatCommas(i.DailyIncomeBeforeTaxStr()),
		models.FormatCommas(i.NetDailyIncomeStr()),
		i.WorkingHours(),
		models.FormatCommas(i.SpecialIncomeRateStr()),
		models.FormatCommas(i.SpecialIncomeBeforeTaxStr()),
		models.FormatCommas(i.NetSpecialIncomeStr()),
		i.GetDeduction(),
		models.FormatCommas(i.TotalVATStr()),
		i.TotalWHTStr(),
		models.FormatCommas(i.TransferAmountStr()),
		i.Note(),
		i.SubmitDateStr(),
	}
	return d
}
