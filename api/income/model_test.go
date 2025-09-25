package income

import (
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestModelIncome(t *testing.T) {
	t.Run("เวลา Add income ควร save ชื่อบัญชี เลขบัญชี และจำนวนเงินด้วย ตอน export จะได้ไม่ต้องคำนวนแล้ว", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.BankAccountName, res.BankAccountName)
		assert.Equal(t, user.BankAccountNumber, res.BankAccountNumber)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, 2000.0, res.DailyRate)
		assert.Equal(t, "38800.00", res.NetDailyIncome)
		assert.Equal(t, "19400.00", res.NetSpecialIncome)
		assert.Equal(t, "58200.00", res.NetIncome)
		assert.Equal(t, "", res.VAT)
		assert.Equal(t, "1800.00", res.WHT)
	})

	t.Run("เวลา Add income ควร save role ด้วย จะได้รู้ว่าเป็น coporate หรือ individual income", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.Role, res.Role)
	})

	t.Run("เวลา Add income ควร save ชื่อ นามสกุล เลขบัตรประชาชนเวลา export ให้บัญชี เค้าจะได้รู้ว่าจ่ายเงินให้ใคร", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "first last", res.Name)
		assert.Equal(t, user.ThaiCitizenID, res.ThaiCitizenID)
	})

	t.Run("เวลา Add income ควร save เบอร์โทรกับ อีเมลด้วยเผื่อตกขบวนเพื่อน ๆ จะได้ช่วยกันตามมากรอกเงินจากหน้า web หน้า individual list ได้", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, user.Phone, res.Phone)
	})

	t.Run("เวลา Add income ควร save วันที่กรอกด้วยจะ เผื่อ export ตอนมีคนตกขบวนจะได้ sort ได้ว่า 2 file รายชื่อต่างกันตรงไหน", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, res.SubmitDate)
	})

	t.Run("เวลา Add income ควร save note ด้วย ไม่รู้ทำไมเหมือนกัน", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncomeReq.Note, res.Note)
	})

	t.Run("เวลา Add income ควร total income ด้วยเพราะ iOS, Andriod และหน้า history ใช้", func(t *testing.T) {
		// ref: https://3.basecamp.com/4877526/buckets/19693649/card_tables/cards/7638832341#__recording_7639315070
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "60000.00", res.TotalIncome)
	})

	t.Run("calculate individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, i.dailyIncome())
		assert.Equal(t, 5*20.0*0.03, i.WitholdingTax(i.dailyIncome()))
		assert.Equal(t, 0.0, i.VAT(i.dailyIncome()))
		assert.Equal(t, 100.0+0-3, i.Net(i.dailyIncome()))
		assert.Equal(t, 10*100.0, i.specialIncome())
		assert.Equal(t, 10*100.0*0.03, i.WitholdingTax(i.specialIncome()))
		assert.Equal(t, 0.0, i.VAT(i.specialIncome()))
		assert.Equal(t, 1000.0+0-30, i.Net(i.specialIncome()))
		assert.Equal(t, i.dailyIncome()+i.specialIncome(), i.totalIncome())
	})

	// begin obsoleted export will be replaced with new export in Aug release
	t.Run("export individual income information เพื่อให้บัญชีติดต่อได้เวลามีปัญหา", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		user.FirstName = "first"
		user.LastName = "last"
		user.ThaiCitizenID = "id"
		user.BankAccountName = "account name"
		user.BankAccountNumber = "0123456789"
		user.Email = "test@example.com"
		req := models.IncomeReq{WorkDate: "20"}
		i := NewIncome(uidFromSession)
		record, _ := i.prepareDataForAddIncome(req, user)
		i = NewIncomeFromRecord(*record)

		csvColumns := i.export(user)

		assert.Equal(t, "first last", csvColumns[NAME_INDEX])
		assert.Equal(t, "id", csvColumns[ID_CARD_INDEX])
		assert.Equal(t, "account name", csvColumns[ACCOUNT_NAME_INDEX])
		assert.Equal(t, `="0123456789"`, csvColumns[ACCOUNT_NUMBER_INDEX])
		assert.Equal(t, "test@example.com", csvColumns[EMAIL_INDEX])
	})

	t.Run("export จำนวนเงินที่ต้องโอนสำหรับ individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)
		record, _ := i.prepareDataForAddIncome(req, user)
		i = NewIncomeFromRecord(*record)
		i.SetLoan(&models.StudentLoan{Amount: 50})

		csvColumns := i.export(user)

		assert.Equal(t, "97.00", csvColumns[NET_DAILY_INCOME_INDEX])
		assert.Equal(t, "970.00", csvColumns[NET_SPECIAL_INCOME_INDEX])
		assert.Equal(t, "50.00", csvColumns[LOAN_DEDUCTION_INDEX])
		assert.Equal(t, "33.00", csvColumns[WITHHOLDING_TAX_INDEX])
		assert.Equal(t, "1,017.00", csvColumns[TRANSFER_AMOUNT_INDEX])
	})

	// end obsoleted export

	t.Run("new export individual income information เพื่อให้บัญชีติดต่อได้เวลามีปัญหา", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		user.FirstName = "first"
		user.LastName = "last"
		user.ThaiCitizenID = "id"
		user.BankAccountName = "account name"
		user.BankAccountNumber = "0123456789"
		user.Email = "test@example.com"
		req := models.IncomeReq{WorkDate: "20"}
		i := NewIncome(uidFromSession)
		record, _ := i.prepareDataForAddIncome(req, user)
		i = NewIncomeFromRecord(*record)

		csvColumns := i.export2()

		assert.Equal(t, "first last", csvColumns[NAME_INDEX])
		assert.Equal(t, "id", csvColumns[ID_CARD_INDEX])
		assert.Equal(t, "account name", csvColumns[ACCOUNT_NAME_INDEX])
		assert.Equal(t, `="0123456789"`, csvColumns[ACCOUNT_NUMBER_INDEX])
		assert.Equal(t, "test@example.com", csvColumns[EMAIL_INDEX])
	})

	t.Run("new export จำนวนเงินที่ต้องโอนสำหรับ individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		dailyIncome := "5"
		workDate := "20"
		specialIncome := "100"
		workingHours := "10"
		u := GivenIndividualUser(uidFromSession, dailyIncome)
		req := models.IncomeReq{
			WorkDate:      workDate,
			SpecialIncome: specialIncome,
			WorkingHours:  workingHours,
		}
		record := CreateIncome(u, req, "note")
		i := NewIncomeFromRecord(*record)
		i.SetLoan(&models.StudentLoan{Amount: 50})

		csvColumns := i.export2()

		assert.Equal(t, "97.00", csvColumns[NET_DAILY_INCOME_INDEX])
		assert.Equal(t, "970.00", csvColumns[NET_SPECIAL_INCOME_INDEX])
		assert.Equal(t, "50.00", csvColumns[LOAN_DEDUCTION_INDEX])
		assert.Equal(t, "33.00", csvColumns[WITHHOLDING_TAX_INDEX])
		assert.Equal(t, "1,017.00", csvColumns[TRANSFER_AMOUNT_INDEX])
	})

	t.Run("calculate individual income โดยไม่ได้กรอก special income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate: "20",
		}
		i := NewIncome(uidFromSession)

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, i.dailyIncome())
		assert.Equal(t, 5*20.0*0.03, i.WitholdingTax(i.dailyIncome()))
		assert.Equal(t, 0.0, i.VAT(i.dailyIncome()))
		assert.Equal(t, 100.0+0-3, i.netDailyIncome())
		assert.Equal(t, "97.00", i.netDailyIncomeStr())
	})

	t.Run("calculate individual special income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{SpecialIncome: "100", WorkingHours: "10"}
		i := NewIncome(uidFromSession)

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 10*100.0, i.specialIncome())
		assert.Equal(t, 10*100.0*0.03, i.WitholdingTax(i.specialIncome()))
		assert.Equal(t, 0.0, i.VAT(i.specialIncome()))
		assert.Equal(t, 1000.0+0-30, i.Net(i.specialIncome()))
		assert.Equal(t, "970.00", i.netSpecialIncomeStr())
	})

	t.Run("calculate individual income สำหรับคนที่มีหนี้ กยศ และบริษัทหักและนำส่งไว้", func(t *testing.T) {
		// เพื่อแก้ปัญหาที่คนไทยหลายคนไม่ยอมใช้หนี้ กยศ ทาง กยศ เลยมีมาตรการให้บริษัท
		// ชำระหนี้ กยศ แทนพนักงาน โดยให้ทางบริษัทหักหนี้ กยศ ออกจากรายได้เลย
		// แต่เพราะชาวออดส์ไม่ใช่พนักงาน คนส่วนใหญ่ก็ยังไปชำระด้วยตัวเอง
		// ยกเว้นบางคนที่ กยศ เข้าใจว่าเป็นพนักงานของเรา ก็จะส่งรายชื่อมาให้หักในเว็บ
		// กยศ ด้านล่าง
		// ref: https://slfrd.dsl.studentloan.or.th/SLFRD/login

		// ใครที่ กยศ ให้หัก เราก็จะหักแล้วไปแจ้งใน basecamp กลุ่ม กยศ ไว้

		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)
		i.SetLoan(&models.StudentLoan{Amount: 50})

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, i.netDailyIncome()+i.netSpecialIncome()-50, i.transferAmount())
	})

	t.Run("หัก ณ ที่จ่าย 3% คิดจากรายได้รวม ไม่นับหนี้ กยศ", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)
		i.SetLoan(&models.StudentLoan{Amount: 50})

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, i.totalIncome()*0.03, i.totalWHT())
	})

	t.Run("calculate corporate income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := models.User{
			ID:   bson.ObjectIdHex(uidFromSession),
			Role: "corporate",
			// ปรกติเวลารายได้เกิน 1.8 ล้าน/ปี ต้องจด VAT
			// ref: https://www.rd.go.th/fileadmin/user_upload/SMEs/infographic/SME_lv1_3.pdf
			Vat:         "Y",
			DailyIncome: "5",
		}
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, i.dailyIncome())
		assert.Equal(t, 5*20.0*0.03, i.WitholdingTax(i.dailyIncome()))
		assert.Equal(t, 7.000000000000001, i.VAT(i.dailyIncome()))
		assert.Equal(t, 100.0+7-3, i.Net(i.dailyIncome()))
		assert.Equal(t, 10*100.0, i.specialIncome())
		assert.Equal(t, 10*100.0*0.03, i.WitholdingTax(i.specialIncome()))
		assert.Equal(t, 10*100.0*0.07, i.VAT(i.specialIncome()))
		assert.Equal(t, 1000.0+70-30, i.Net(i.specialIncome()))
	})

	t.Run("test export to SAP transaction should format correctly", func(t *testing.T) {
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		i := incomeMock.MockSoloCorporateIncome

		txn, _ := NewIncomeFromRecord(i).exportSAP(dateEff)

		assert.Equal(t, "TXN", txn[SAP_TXN_INDEX])
		assert.Equal(t, "บจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           ", txn[SAP_PAYER_NAME_INDEX])
		assert.Equal(t, "บจก. โซโล่ เลเวลลิ่ง                                                                                                              ", txn[SAP_PAYEE_NAME_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_MALE_TO_NAME_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY1_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY2_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY3_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY4_INDEX])
		assert.Equal(t, "          ", txn[SAP_ZIPCODE_INDEX])
		assert.Equal(t, "                ", txn[SAP_CUSTOMER_REF_INDEX])
		assert.Equal(t, "29092025", txn[SAP_DATE_EFFECTIVE_INDEX])
		assert.Equal(t, "29092025", txn[SAP_DATE_PICKUP_INDEX])
		assert.Equal(t, "THB", txn[SAP_CURRENCY_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_EMPTY_1_INDEX])
		assert.Equal(t, "00011595873         ", txn[SAP_COMPANY_ACCNO_INDEX])
		assert.Equal(t, "000000054704.00", txn[SAP_AMOUNT_INDEX])
		assert.Equal(t, "0110246         ", txn[SAP_PAYEE_BANK_CODE_INDEX])
		assert.Equal(t, "02462737202         ", txn[SAP_ACCOUNTNO_INDEX])
		assert.Equal(t, "0400", txn[SAP_UNKNOW_1_INDEX])
		assert.Equal(t, "  ", txn[SAP_EMPTY_2_INDEX])
		assert.Equal(t, "                    ", txn[SAP_EMPTY_3_INDEX])
		assert.Equal(t, "     ", txn[SAP_ADVICEMODE2_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_FAXNO_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_EMAIL_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_SMSNO_INDEX])
		assert.Equal(t, "OUR          ", txn[SAP_CHARGE_ON_INDEX])
		assert.Equal(t, "DCR", txn[SAP_PRODUCT_INDEX])
		assert.Equal(t, "     ", txn[SAP_SCHEDULE_INDEX])
		assert.Equal(t, "                                  ", txn[SAP_EMPTY_4_INDEX])
		assert.Equal(t, "                                                                                                         ", txn[SAP_DOCREQ_INDEX])
		assert.Equal(t, "                                                                                                                                                                                                                                                                                                       ", txn[SAP_EMPTY_5_INDEX])
		assert.Equal(t, "END", txn[SAP_END_INDEX])

	})

	t.Run("test export to SAP wht should format correctly", func(t *testing.T) {
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		i := incomeMock.MockSoloCorporateIncome

		_, wht := NewIncomeFromRecord(i).exportSAP(dateEff)

		assert.Equal(t, "WHT", wht[SAP_WHT_WHT_INDEX])
		assert.Equal(t, "             ", wht[SAP_WHT_EMPTY_1_INDEX])
		assert.Equal(t, "0105556110718", wht[SAP_WHT_TAX_ID_INDEX])
		assert.Equal(t, "  ", wht[SAP_WHT_EMPTY_3_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_4_INDEX])
		assert.Equal(t, "  ", wht[SAP_WHT_EMPTY_5_INDEX])
		assert.Equal(t, "                                   ", wht[SAP_WHT_EMPTY_6_INDEX])
		assert.Equal(t, "     ", wht[SAP_WHT_EMPTY_7_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_8_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_9_INDEX])
		assert.Equal(t, "  ", wht[SAP_WHT_EMPTY_10_INDEX])
		assert.Equal(t, "                                   ", wht[SAP_WHT_EMPTY_11_INDEX])
		assert.Equal(t, "     ", wht[SAP_WHT_EMPTY_12_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_13_INDEX])
		assert.Equal(t, "                                                                                                                                                ", wht[SAP_WHT_EMPTY_14_INDEX])
		assert.Equal(t, "บจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           ", wht[SAP_WHT_COM_NAME])
		assert.Equal(t, "2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900                                                                                                                ", wht[SAP_WHT_ADDRESS])
		assert.Equal(t, "                                                                                                                        ", wht[SAP_WHT_EMPTY_17])
		assert.Equal(t, "                                                                                                                                                                ", wht[SAP_WHT_EMPTY_18])
		assert.Equal(t, "                    ", wht[SAP_WHT_EMPTY_19])
		assert.Equal(t, AddBlank("", 938), wht[SAP_WHT_EMPTY_20])

	})

}

func TestModelIncomes(t *testing.T) {
	t.Run("test export to CSV when there is 0 income", func(t *testing.T) {
		records := []*models.Income{}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.toCSV()

		assert.NotNil(t, csv)
		headerLength := 1
		assert.Equal(t, headerLength, len(csv))
	})

	t.Run("test export to CSV when there is 1 income", func(t *testing.T) {
		records := []*models.Income{
			{ID: "incomeId"},
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		csv, _ := incomes.toCSV()

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

		csv, _ := incomes.toCSV()

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

		csv, _ := incomes.toCSV()

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

		csv, _ := incomes.toCSV()

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

		_, updatedIncomeIds := incomes.toCSV()

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

		_, updatedIncomeIds := incomes.toCSV()

		assert.NotNil(t, updatedIncomeIds)
		assert.Equal(t, 1, len(updatedIncomeIds))
	})

	t.Run("test toCSVasSAP", func(t *testing.T) {

		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)

		records := []*models.Income{
			&incomeMock.MockSoloCorporateIncome,
			&incomeMock.MockSwardCorporateIncome,
		}
		incomes := NewIncomes(records, models.StudentLoanList{})

		i, _ := incomes.toCSVasSAP(dateEff)

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
