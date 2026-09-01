package file

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestPeakHeaders(t *testing.T) {
	actual := peakHeaders()
	expected := []string{
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
	assert.Equal(t, expected, actual)
}

func TestToPeakCSVMatchesSample(t *testing.T) {
	period := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	documentDate := time.Date(2026, 7, 25, 0, 0, 0, 0, time.UTC)

	incomes := []*models.Income{
		{
			UserID:          "user-vat",
			TotalIncome:     "150000.00",
			IsVATRegistered: true,
			WHTRate:         0.05,
			PeakCode:        "C00145",
		},
		{
			UserID:          "user-no-vat",
			TotalIncome:     "40000",
			IsVATRegistered: false,
			WHTRate:         0.05,
			PeakCode:        "C00464",
		},
	}

	rows := ToPeakCSV(incomes, period, documentDate)
	assert.Equal(t, 3, len(rows))

	vatRow := rows[1]
	assert.Equal(t, "1", vatRow[peakSeqIndex])
	assert.Equal(t, "20260725", vatRow[peakDocumentDateIndex])
	assert.Equal(t, "", vatRow[peakReferenceIndex])
	assert.Equal(t, "C00145", vatRow[peakContactIndex])
	assert.Equal(t, "", vatRow[peakTaxIDIndex])
	assert.Equal(t, "1", vatRow[peakPriceTypeIndex])
	assert.Equal(t, "510111", vatRow[peakAccountIndex])
	assert.Equal(t, "ค่าพัฒนาและสอนโปรแกรม(บุคคล) 07.2026", vatRow[peakDescriptionIndex])
	assert.Equal(t, "1", vatRow[peakQuantityIndex])
	assert.Equal(t, "150000", vatRow[peakUnitPriceIndex])
	assert.Equal(t, "7%", vatRow[peakTaxRateIndex])
	assert.Equal(t, "5%", vatRow[peakWHTIndex])
	assert.Equal(t, "BCA001", vatRow[peakPaymentIndex])
	assert.Equal(t, "150000", vatRow[peakPaidAmountIndex])
	assert.Equal(t, "1", vatRow[peakPNDIndex])
	assert.Equal(t, "", vatRow[peakNoteIndex])
	assert.Equal(t, "", vatRow[peakCategoryIndex])

	noVATRow := rows[2]
	assert.Equal(t, "2", noVATRow[peakSeqIndex])
	assert.Equal(t, "C00464", noVATRow[peakContactIndex])
	assert.Equal(t, "", noVATRow[peakTaxIDIndex])
	assert.Equal(t, "40000", noVATRow[peakUnitPriceIndex])
	assert.Equal(t, "NO", noVATRow[peakTaxRateIndex])
	assert.Equal(t, "40000", noVATRow[peakPaidAmountIndex])
}

func TestToPeakCSVFallsBackToThaiCitizenIDWhenPeakCodeMissing(t *testing.T) {
	period := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	documentDate := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)
	incomes := []*models.Income{
		{
			UserID:        "user-no-peak",
			TotalIncome:   "10000.00",
			ThaiCitizenID: "1234567890123",
		},
	}

	rows := ToPeakCSV(incomes, period, documentDate)
	assert.Equal(t, "", rows[1][peakContactIndex])
	assert.Equal(t, "1234567890123", rows[1][peakTaxIDIndex])
	assert.Equal(t, "ค่าพัฒนาและสอนโปรแกรม(บุคคล) 08.2026", rows[1][peakDescriptionIndex])
	assert.Equal(t, "20260815", rows[1][peakDocumentDateIndex])
}

func TestToPeakCSVUsesWHTRateSavedOnIncome(t *testing.T) {
	period := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	documentDate := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)
	incomes := []*models.Income{
		{UserID: "old-month", TotalIncome: "10000", WHTRate: 0.03},
		{UserID: "new-month", TotalIncome: "10000", WHTRate: 0.05},
		{UserID: "legacy", TotalIncome: "10000"},
	}

	rows := ToPeakCSV(incomes, period, documentDate)

	assert.Equal(t, "3%", rows[1][peakWHTIndex])
	assert.Equal(t, "5%", rows[2][peakWHTIndex])
	assert.Equal(t, "3%", rows[3][peakWHTIndex])
}

func TestPeakDocumentDateUsesThailandOffset(t *testing.T) {
	// 17:00 UTC + 7h = next calendar day in Thailand
	documentDate := time.Date(2026, 8, 15, 17, 0, 0, 0, time.UTC)
	assert.Equal(t, "20260816", peakDocumentDate(documentDate))
}

func TestToPeakCSVHeaderOnlyWhenNoIncomes(t *testing.T) {
	period := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	rows := ToPeakCSV(nil, period, period)
	assert.Equal(t, 1, len(rows))
	assert.Equal(t, peakHeaders(), rows[0])
}

func TestPeakAmountStripsCommasAndTrailingZeroDecimals(t *testing.T) {
	assert.Equal(t, "150000", peakAmount("150,000.00"))
	assert.Equal(t, "150000.50", peakAmount("150000.50"))
}
