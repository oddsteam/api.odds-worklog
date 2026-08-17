package file

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

const (
	peakAccountCode    = "510111"
	peakPaymentChannel = "BCA001"
	peakPriceType      = "1"
	peakQuantity       = "1"
	peakWHT            = "3%"
	peakPND            = "1"
	peakVATRate        = "7%"
	peakNoVAT          = "NO"
	thailandOffset     = 7 * time.Hour
)

const (
	peakSeqIndex = iota
	peakDocumentDateIndex
	peakReferenceIndex
	peakContactIndex
	peakTaxIDIndex
	peakBranchIndex
	peakInvoiceNoIndex
	peakInvoiceDateIndex
	peakVATRecordDateIndex
	peakPriceTypeIndex
	peakAccountIndex
	peakDescriptionIndex
	peakQuantityIndex
	peakUnitPriceIndex
	peakTaxRateIndex
	peakWHTIndex
	peakPaymentIndex
	peakPaidAmountIndex
	peakPNDIndex
	peakNoteIndex
	peakCategoryIndex
)

type peakCSVWriter struct{}

func NewPeakCSVWriter() *peakCSVWriter {
	return &peakCSVWriter{}
}

func (w *peakCSVWriter) WriteFile(name string, incomes []*models.Income, peakCodeByUserID map[string]string, period, documentDate time.Time) (string, error) {
	rows := ToPeakCSV(incomes, peakCodeByUserID, period, documentDate)
	if len(rows) == 1 {
		return "", errors.New("no data for export to CSV file")
	}

	file, filename, err := CreateFile(name)
	if err != nil {
		return "", err
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(rows)
	csvWriter.Flush()
	defer file.Close()
	return filename, nil
}

func ToPeakCSV(incomes []*models.Income, peakCodeByUserID map[string]string, period, documentDate time.Time) [][]string {
	rows := [][]string{peakHeaders()}
	seq := 0
	for _, income := range incomes {
		if income == nil {
			continue
		}
		seq++
		peakCode := ""
		if peakCodeByUserID != nil {
			peakCode = peakCodeByUserID[income.UserID]
		}
		rows = append(rows, peakRow(seq, income, peakCode, period, documentDate))
	}
	return rows
}

func peakHeaders() []string {
	return []string{
		"ลำดับที่*",
		"วันที่เอกสาร",
		"อ้างอิงถึง",
		"ผู้รับเงิน/คู่ค้า",
		"เลขทะเบียน 13 หลัก",
		"เลขสาขา 5 หลัก",
		"เลขที่ใบกำกับฯ (ถ้ามี)",
		"วันที่ใบกำกับฯ (ถ้ามี)",
		"วันที่บันทึกภาษีซื้อ (ถ้ามี)",
		"ประเภทราคา",
		"บัญชี",
		"คำอธิบาย",
		"จำนวน",
		"ราคาต่อหน่วย",
		"อัตราภาษี",
		"หัก ณ ที่จ่าย (ถ้ามี)",
		"ชำระโดย",
		"จำนวนเงินที่ชำระ",
		"ภ.ง.ด. (ถ้ามี)",
		"หมายเหตุ",
		"กลุ่มจัดประเภท",
	}
}

func peakRow(seq int, income *models.Income, peakCode string, period, documentDate time.Time) []string {
	contact := peakCode
	taxID := ""
	if contact == "" {
		taxID = income.ThaiCitizenID
	}
	amount := peakAmount(income.TotalIncome)
	return []string{
		strconv.Itoa(seq),
		peakDocumentDate(documentDate),
		"",
		contact,
		taxID,
		"",
		"",
		"",
		"",
		peakPriceType,
		peakAccountCode,
		peakDescription(period),
		peakQuantity,
		amount,
		peakTaxRate(income.IsVATRegistered),
		peakWHT,
		peakPaymentChannel,
		amount,
		peakPND,
		"",
		"",
	}
}

func peakDescription(period time.Time) string {
	return fmt.Sprintf("ค่าพัฒนาและสอนโปรแกรม(บุคคล) %02d.%d", int(period.Month()), period.Year())
}

func peakDocumentDate(t time.Time) string {
	return t.Add(thailandOffset).Format("20060102")
}

func peakTaxRate(vatRegistered bool) string {
	if vatRegistered {
		return peakVATRate
	}
	return peakNoVAT
}

func peakAmount(totalIncome string) string {
	s := strings.ReplaceAll(totalIncome, ",", "")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	if f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}
