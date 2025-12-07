package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/entity"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestCSVWriter(t *testing.T) {
	t.Run("test updatedIncomeIds when มีคนตกขบวน", func(t *testing.T) {
		users := []*models.User{
			{ID: "id1"},
			{ID: "id2"},
		}
		records := []*models.Income{
			{ID: "incomeId1", UserID: users[0].ID.Hex()},
		}
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		_, updatedIncomeIds := ToCSV(*incomes)

		assert.NotNil(t, updatedIncomeIds)
		assert.Equal(t, 1, len(updatedIncomeIds))
	})

	t.Run("export individual income information เพื่อให้บัญชีติดต่อได้เวลามีปัญหา", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := entity.GivenIndividualUser(uidFromSession, "5")
		user.FirstName = "first"
		user.LastName = "last"
		user.ThaiCitizenID = "id"
		user.BankAccountName = "account name"
		user.BankAccountNumber = "0123456789"
		user.Email = "test@example.com"
		req := entity.IncomeReq{WorkDate: "20"}
		record := entity.CreateIncome(user, req, "note")
		i := entity.NewIncomeFromRecord(*record)

		csvColumns := export(*i)

		assert.Equal(t, "first last", csvColumns[NAME_INDEX])
		assert.Equal(t, "id", csvColumns[ID_CARD_INDEX])
		assert.Equal(t, "account name", csvColumns[ACCOUNT_NAME_INDEX])
		assert.Equal(t, `="0123456789"`, csvColumns[ACCOUNT_NUMBER_INDEX])
		assert.Equal(t, "test@example.com", csvColumns[EMAIL_INDEX])
	})

	t.Run("export จำนวนเงินที่ต้องโอนสำหรับ individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		dailyIncome := "5"
		workDate := "20"
		specialIncome := "100"
		workingHours := "10"
		u := entity.GivenIndividualUser(uidFromSession, dailyIncome)
		req := entity.IncomeReq{
			WorkDate:      workDate,
			SpecialIncome: specialIncome,
			WorkingHours:  workingHours,
		}
		record := entity.CreateIncome(u, req, "note")
		i := entity.NewIncomeFromRecord(*record)
		i.SetLoan(&models.StudentLoan{})

		csvColumns := export(*i)

		assert.Equal(t, "97.00", csvColumns[NET_DAILY_INCOME_INDEX])
		assert.Equal(t, "970.00", csvColumns[NET_SPECIAL_INCOME_INDEX])
		assert.Equal(t, "0.00", csvColumns[LOAN_DEDUCTION_INDEX])
		assert.Equal(t, "33.00", csvColumns[WITHHOLDING_TAX_INDEX])
		assert.Equal(t, "1,067.00", csvColumns[TRANSFER_AMOUNT_INDEX])
		assert.Equal(t, "note", csvColumns[NOTE_INDEX])
	})

}

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
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		csv, _ := ToCSV(*incomes)

		assert.NotNil(t, csv)
		headerLength := 1
		assert.Equal(t, headerLength, len(csv))
	})

	t.Run("test export to CSV when there is 1 income", func(t *testing.T) {
		records := []*models.Income{
			{ID: "incomeId"},
		}
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		csv, _ := ToCSV(*incomes)

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
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		csv, _ := ToCSV(*incomes)

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
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		csv, _ := ToCSV(*incomes)

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
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		csv, _ := ToCSV(*incomes)

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
		incomes := entity.NewIncomes(records, models.StudentLoanList{})

		_, updatedIncomeIds := ToCSV(*incomes)

		assert.NotNil(t, updatedIncomeIds)
		assert.Equal(t, 1, len(updatedIncomeIds))
	})
}
