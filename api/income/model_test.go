package income

import (
	"testing"

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
		assert.Equal(t, "40000.00", res.TotalIncome)
	})

	t.Run("calculate individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, i.totalIncome())
		assert.Equal(t, 5*20.0*0.03, i.WitholdingTax(i.totalIncome()))
		assert.Equal(t, 0.0, i.VAT(i.totalIncome()))
		assert.Equal(t, 100.0+0-3, i.Net(i.totalIncome()))
		assert.Equal(t, 10*100.0, i.specialIncome())
		assert.Equal(t, 10*100.0*0.03, i.WitholdingTax(i.specialIncome()))
		assert.Equal(t, 0.0, i.VAT(i.specialIncome()))
		assert.Equal(t, 1000.0+0-30, i.Net(i.specialIncome()))
	})

	// begin obsoleted export will be replaced with new export in Aug release
	t.Run("export individual income information เพื่อให้บัญชีติดต่อได้เวลามีปัญหา", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
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

		assert.Equal(t, "first last", csvColumns[0])
		assert.Equal(t, "id", csvColumns[1])
		assert.Equal(t, "account name", csvColumns[2])
		assert.Equal(t, `="0123456789"`, csvColumns[3])
		assert.Equal(t, "test@example.com", csvColumns[4])
	})

	t.Run("export จำนวนเงินที่ต้องโอนสำหรับ individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
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

		assert.Equal(t, "97.00", csvColumns[5])
		assert.Equal(t, "970.00", csvColumns[6])
		assert.Equal(t, "50.00", csvColumns[7])
		assert.Equal(t, "33.00", csvColumns[8])
		assert.Equal(t, "1,017.00", csvColumns[9])
	})

	// end obsoleted export

	t.Run("new export individual income information เพื่อให้บัญชีติดต่อได้เวลามีปัญหา", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
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

		assert.Equal(t, "first last", csvColumns[0])
		assert.Equal(t, "id", csvColumns[1])
		assert.Equal(t, "account name", csvColumns[2])
		assert.Equal(t, `="0123456789"`, csvColumns[3])
		assert.Equal(t, "test@example.com", csvColumns[4])
	})

	t.Run("new export จำนวนเงินที่ต้องโอนสำหรับ individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		i := NewIncome(uidFromSession)
		record, _ := i.prepareDataForAddIncome(req, user)
		i = NewIncomeFromRecord(*record)
		i.SetLoan(&models.StudentLoan{Amount: 50})

		csvColumns := i.export2()

		assert.Equal(t, "97.00", csvColumns[5])
		assert.Equal(t, "970.00", csvColumns[6])
		assert.Equal(t, "33.00", csvColumns[8])
		assert.Equal(t, "50.00", csvColumns[7])
		assert.Equal(t, "1,017.00", csvColumns[9])
	})

	t.Run("calculate individual income โดยไม่ได้กรอก special income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
		req := models.IncomeReq{
			WorkDate: "20",
		}
		i := NewIncome(uidFromSession)

		err := i.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, i.totalIncome())
		assert.Equal(t, 5*20.0*0.03, i.WitholdingTax(i.totalIncome()))
		assert.Equal(t, 0.0, i.VAT(i.totalIncome()))
		assert.Equal(t, 100.0+0-3, i.netDailyIncome())
		assert.Equal(t, "97.00", i.netDailyIncomeStr())
	})

	t.Run("calculate individual special income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := givenIndividualUser(uidFromSession, "5")
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
		user := givenIndividualUser(uidFromSession, "5")
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
		assert.Equal(t, 5*20.0, i.totalIncome())
		assert.Equal(t, 5*20.0*0.03, i.WitholdingTax(i.totalIncome()))
		assert.Equal(t, 7.000000000000001, i.VAT(i.totalIncome()))
		assert.Equal(t, 100.0+7-3, i.Net(i.totalIncome()))
		assert.Equal(t, 10*100.0, i.specialIncome())
		assert.Equal(t, 10*100.0*0.03, i.WitholdingTax(i.specialIncome()))
		assert.Equal(t, 10*100.0*0.07, i.VAT(i.specialIncome()))
		assert.Equal(t, 1000.0+70-30, i.Net(i.specialIncome()))
	})
}

func givenIndividualUser(uidFromSession string, dailyIncome string) models.User {
	return models.User{
		ID:          bson.ObjectIdHex(uidFromSession),
		Role:        "individual",
		Vat:         "N",
		DailyIncome: dailyIncome,
	}
}

func TestModelIncomes(t *testing.T) {
	t.Run("test export to CSV when there is 0 income", func(t *testing.T) {
		users := []*models.User{{ID: "id"}}
		records := []*models.Income{}
		incomes := NewIncomes(records, models.StudentLoanList{}, users)

		csv, _ := incomes.toCSV()

		assert.NotNil(t, csv)
		headerLenght := 1
		assert.Equal(t, headerLenght, len(csv))
	})

	t.Run("test export to CSV when there is 1 income", func(t *testing.T) {
		users := []*models.User{
			{ID: "id"},
		}
		records := []*models.Income{
			{ID: "incomeId", UserID: users[0].ID.Hex()},
		}
		incomes := NewIncomes(records, models.StudentLoanList{}, users)

		csv, _ := incomes.toCSV()

		assert.NotNil(t, csv)
		headerLenght := 1
		incomeCount := 1
		assert.Equal(t, headerLenght+incomeCount, len(csv))
	})

	t.Run("test export to CSV when there is n incomes", func(t *testing.T) {
		users := []*models.User{
			{ID: "id1"},
			{ID: "id2"},
		}
		records := []*models.Income{
			{ID: "incomeId1", UserID: users[0].ID.Hex()},
			{ID: "incomeId2", UserID: users[1].ID.Hex()},
		}
		incomes := NewIncomes(records, models.StudentLoanList{}, users)

		csv, _ := incomes.toCSV()

		assert.NotNil(t, csv)
		headerLenght := 1
		incomeCount := 2
		assert.Equal(t, headerLenght+incomeCount, len(csv))
	})

	t.Run("test export to CSV when มีคนตกขบวน", func(t *testing.T) {
		users := []*models.User{
			{ID: "id1"},
			{ID: "id2"},
		}
		records := []*models.Income{
			{ID: "incomeId1", UserID: users[0].ID.Hex()},
		}
		incomes := NewIncomes(records, models.StudentLoanList{}, users)

		csv, _ := incomes.toCSV()

		assert.NotNil(t, csv)
		headerLenght := 1
		incomeCount := 1
		assert.Equal(t, headerLenght+incomeCount, len(csv))
	})

	t.Run("test should also return updatedIncomeIds", func(t *testing.T) {
		// เราจะได้ mark ว่า income เหล่านี้ถูก export ออกไปแล้ว
		// เวลา export different individuals (คนที่ยังไม่ถูก export, เพราะตัดรอบไปก่อน)
		// จะได้รู้ว่าใครบ้างที่ export ไปแล้ว
		users := []*models.User{
			{ID: "id"},
		}
		records := []*models.Income{
			{ID: "incomeId", UserID: users[0].ID.Hex()},
		}
		incomes := NewIncomes(records, models.StudentLoanList{}, users)

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
		incomes := NewIncomes(records, models.StudentLoanList{}, users)

		_, updatedIncomeIds := incomes.toCSV()

		assert.NotNil(t, updatedIncomeIds)
		assert.Equal(t, 1, len(updatedIncomeIds))
	})
}
