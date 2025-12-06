package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestCSVHeaders(t *testing.T) {
	actual := createHeaders()
	expected := [...]string{"Vendor Code", "ชื่อบัญชี", "Payment method", "เลขบัญชี", "ชื่อ", "เลขบัตรประชาชน",
		"อีเมล", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ",
		"หัก ณ ที่จ่าย", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก",
	}
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

func TestModelIncomes(t *testing.T) {
	t.Run("test export to CSV when there is 0 income", func(t *testing.T) {
		records := []*models.Income{}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.ToCSV()

		assert.NotNil(t, csv)
		headerLength := 1
		assert.Equal(t, headerLength, len(csv))
	})

	t.Run("test export to CSV when there is 1 income", func(t *testing.T) {
		records := []*models.Income{
			{ID: "incomeId"},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.ToCSV()

		assert.NotNil(t, csv)
		headerLength := 1
		incomeCount := 1
		assert.Equal(t, headerLength+incomeCount, len(csv))
	})

	t.Run("test export to CSV when there is n incomes", func(t *testing.T) {
		records := []*models.Income{
			{ID: "incomeId1"},
			{ID: "incomeId2"},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.ToCSV()

		assert.NotNil(t, csv)
		headerLength := 1
		incomeCount := 2
		assert.Equal(t, headerLength+incomeCount, len(csv))
	})

	t.Run("test export to CSV with running vendor codes", func(t *testing.T) {
		records := []*models.Income{
			{ID: "incomeId1"},
			{ID: "incomeId2"},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.ToCSV()

		assert.NotNil(t, csv)
		headerLength := 1
		incomeCount := 2
		assert.Equal(t, headerLength+incomeCount, len(csv))
		assert.Equal(t, "AAA", csv[1][VENDOR_CODE_INDEX])
		assert.Equal(t, "AAB", csv[2][VENDOR_CODE_INDEX])
	})

	t.Run("test export to CSV when มีคนตกขบวน", func(t *testing.T) {
		users := []*models.User{
			{ID: "id1"},
			{ID: "id2"},
		}
		records := []*models.Income{
			{ID: "incomeId1", UserID: users[0].ID.Hex()},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.ToCSV()

		assert.NotNil(t, csv)
		headerLength := 1
		incomeCount := 1
		assert.Equal(t, headerLength+incomeCount, len(csv))
	})

	t.Run("test should also return updatedIncomeIds", func(t *testing.T) {
		// เราจะได้ mark ว่า income เหล่านี้ถูก export ออกไปแล้ว
		// เวลา export different individuals (คนที่ยังไม่ถูก export, เพราะตัดรอบไปก่อน)
		// จะได้รู้ว่าใครบ้างที่ export ไปแล้ว
		records := []*models.Income{
			{ID: "incomeId"},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		_, updatedIncomeIds := incomes.ToCSV()

		assert.NotNil(t, updatedIncomeIds)
		assert.Equal(t, 1, len(updatedIncomeIds))
	})

	t.Run("test updatedIncomeIds when มีคนตกขบวน", func(t *testing.T) {
		users := []*models.User{
			{ID: "id1"},
			{ID: "id2"},
		}
		records := []*models.Income{
			{ID: "incomeId1", UserID: users[0].ID.Hex()},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		_, updatedIncomeIds := incomes.ToCSV()

		assert.NotNil(t, updatedIncomeIds)
		assert.Equal(t, 1, len(updatedIncomeIds))
	})

	t.Run("test transform to SAP format", func(t *testing.T) {

		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)

		records := []*models.Income{
			&MockSoloCorporateIncome,
			&MockSwardCorporateIncome,
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		i, _ := incomes.ToSAP(dateEff)

		assert.Equal(t, "TXNบจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           บจก. โซโล่ เลเวลลิ่ง                                                                                                                                                                                                                                                                                                                                                2909202529092025THB                                                  00011595873         000000054704.000110246         02462737202         0400                                                                                                                                                                                 OUR          DCR                                                                                                                                                                                                                                                                                                                                                                                                                                                       END", strings.Join(i[0], ""))
		assert.Equal(t, "WHT             0105556110718  000000000000.00                                          000000000000.00000000000000.00                                          000000000000.00                                                                                                                                                บจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      ", strings.Join(i[1], ""))
		assert.Equal(t, "TXNบจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           บจ. ดาบพิฆาตอสูร                                                                                                                                                                                                                                                                                                                                                    2909202529092025THB                                                  00011595873         000000005470.400110110         01102480447         0400                                                                                                                                                                                 OUR          DCR                                                                                                                                                                                                                                                                                                                                                                                                                                                       END", strings.Join(i[2], ""))
		assert.Equal(t, "WHT             0105556110718  000000000000.00                                          000000000000.00000000000000.00                                          000000000000.00                                                                                                                                                บจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      ", strings.Join(i[3], ""))

	})
}

func TestModelVendorCode(t *testing.T) {
	t.Run("test large index vendor code", func(t *testing.T) {
		vc := VendorCode{index: 381}
		assert.Equal(t, "AOR", vc.String())
	})
}
