package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

// Local test mock for IndividualUser1
var testIndividualUser1 = User{
	ID:                bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39531"),
	Role:              "individual",
	FirstName:         "first",
	LastName:          "last",
	Email:             "email@example.com",
	BankAccountName:   "ชื่อ นามสกุล",
	BankAccountNumber: "0531231231",
	ThaiCitizenID:     "ThaiCitizenID",
	Vat:               "N",
	DailyIncome:       "2000",
	StatusTavi:        true,
	Address:           "address",
	StartDate:         "2022-01-01",
}

func TestPayroll(t *testing.T) {
	t.Run("เวลา Add income ควร save ชื่อบัญชี เลขบัญชี และจำนวนเงินด้วย ตอน export จะได้ไม่ต้องคำนวนแล้ว", func(t *testing.T) {
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.BankAccountName, res.BankAccountName)
		assert.Equal(t, user.BankAccountNumber, res.BankAccountNumber)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, 2000.0, res.DailyRate)
		assert.Equal(t, "40000.00", res.DailyIncomeBeforeTax)
		assert.Equal(t, "20000.00", res.SpecialIncomeBeforeTax)
		assert.Equal(t, "38000.00", res.NetDailyIncome)
		assert.Equal(t, "19000.00", res.NetSpecialIncome)
		assert.Equal(t, "57000.00", res.NetIncome)
		assert.Equal(t, "", res.VAT)
		assert.Equal(t, "3000.00", res.WHT)
		assert.Equal(t, 0.05, res.WHTRate)
	})

	t.Run("เวลา Add income ควร save role ด้วย จะได้รู้ว่าเป็น coporate หรือ individual income", func(t *testing.T) {
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.Role, res.Role)
	})

	t.Run("เวลา Add income ควร save ชื่อ นามสกุล เลขบัตรประชาชนเวลา export ให้บัญชี เค้าจะได้รู้ว่าจ่ายเงินให้ใคร", func(t *testing.T) {
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "first last", res.Name)
		assert.Equal(t, user.ThaiCitizenID, res.ThaiCitizenID)
	})

	t.Run("เวลา Add income ควร save เบอร์โทรกับ อีเมลด้วยเผื่อตกขบวนเพื่อน ๆ จะได้ช่วยกันตามมากรอกเงินจากหน้า web หน้า individual list ได้", func(t *testing.T) {
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, user.Phone, res.Phone)
	})

	t.Run("เวลา Add income ควร save วันที่กรอกด้วยจะ เผื่อ export ตอนมีคนตกขบวนจะได้ sort ได้ว่า 2 file รายชื่อต่างกันตรงไหน", func(t *testing.T) {
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, res.SubmitDate)
	})

	t.Run("เวลา Add income ควร save note ด้วย ไม่รู้ทำไมเหมือนกัน", func(t *testing.T) {
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, MockIncomeReq.Note, res.Note)
	})

	t.Run("เวลา Add income ควร total income ด้วยเพราะ iOS, Andriod และหน้า history ใช้", func(t *testing.T) {
		// ref: https://3.basecamp.com/4877526/buckets/19693649/card_tables/cards/7638832341#__recording_7639315070
		user := testIndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)

		res, err := p.prepareDataForAddIncome(MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "60000.00", res.TotalIncome)
	})

	t.Run("calculate individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		p := NewPayroll(uidFromSession)

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, p.dailyIncome())
		assert.Equal(t, 5*20.0*0.05, p.WitholdingTax(p.dailyIncome()))
		assert.Equal(t, 0.0, p.VAT(p.dailyIncome()))
		assert.Equal(t, 100.0+0-5, p.Net(p.dailyIncome()))
		assert.Equal(t, 10*100.0, p.specialIncome())
		assert.Equal(t, 10*100.0*0.05, p.WitholdingTax(p.specialIncome()))
		assert.Equal(t, 0.0, p.VAT(p.specialIncome()))
		assert.Equal(t, 1000.0+0-50, p.Net(p.specialIncome()))
		assert.Equal(t, p.dailyIncome()+p.specialIncome(), p.totalIncome())
	})

	t.Run("calculate individual income โดยไม่ได้กรอก special income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := IncomeReq{
			WorkDate: "20",
		}
		p := NewPayroll(uidFromSession)

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, p.dailyIncome())
		assert.Equal(t, 5*20.0*0.05, p.WitholdingTax(p.dailyIncome()))
		assert.Equal(t, 0.0, p.VAT(p.dailyIncome()))
		assert.Equal(t, 100.0+0-5, p.netDailyIncome())
		assert.Equal(t, "95.00", p.NetDailyIncomeStr())
	})

	t.Run("calculate individual special income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := IncomeReq{SpecialIncome: "100", WorkingHours: "10"}
		p := NewPayroll(uidFromSession)

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 10*100.0, p.specialIncome())
		assert.Equal(t, 10*100.0*0.05, p.WitholdingTax(p.specialIncome()))
		assert.Equal(t, 0.0, p.VAT(p.specialIncome()))
		assert.Equal(t, 1000.0+0-50, p.Net(p.specialIncome()))
		assert.Equal(t, "950.00", p.NetSpecialIncomeStr())
	})

	t.Run("calculate decimal daily and special income and apply 5 percent withholding tax to both", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "100")
		req := IncomeReq{
			WorkDate:      "20.50",
			SpecialIncome: "100.50",
			WorkingHours:  "10",
		}
		p := NewPayroll(uidFromSession)

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 2050.0, p.dailyIncome())
		assert.Equal(t, 102.5, p.WitholdingTax(p.dailyIncome()))
		assert.Equal(t, "1947.50", p.NetDailyIncomeStr())
		assert.Equal(t, 1005.0, p.specialIncome())
		assert.Equal(t, 50.25, p.WitholdingTax(p.specialIncome()))
		assert.Equal(t, "954.75", p.NetSpecialIncomeStr())
		assert.Equal(t, "152.75", p.TotalWHTStr())
		assert.Equal(t, "3055.00", p.totalIncomeStr())
		assert.Equal(t, "2902.25", p.TransferAmountStr())
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
		req := IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		p := NewPayroll(uidFromSession)
		p.SetLoan(&StudentLoan{Amount: 50})

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, p.netDailyIncome()+p.netSpecialIncome()-50, p.TransferAmount())
	})

	t.Run("หัก ณ ที่จ่าย 5% คิดจากรายได้รวม ไม่นับหนี้ กยศ", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := GivenIndividualUser(uidFromSession, "5")
		req := IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		p := NewPayroll(uidFromSession)
		p.SetLoan(&StudentLoan{Amount: 50})

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, p.totalIncome()*0.05, p.totalWHT())
	})

	t.Run("student loan is used as deduction for foreign student who does not require social security", func(t *testing.T) {
		// นักศึกษาต่างด้าวที่ยังไม่บรรจุเป็นพนักงานประจำ จะไม่มีประกันสังคม จึงไม่ต้อง
		// หักประกันสังคม 270 บาท เหมือนคนไทย เราใส่ช่อง deduction เป็นลบ 270
		// บาท เพื่อคืนเงินที่หักประกันสังคมคืนไป
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		p := NewPayroll(uidFromSession)
		p.SetLoan(&StudentLoan{Amount: -270})
		user := GivenIndividualUser(uidFromSession, "5")
		req := IncomeReq{
			SpecialIncome: "100",
			WorkingHours:  "10",
		}

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, p.netSpecialIncome()+270, p.TransferAmount())
	})

	t.Run("เวลา Update income ควรเก็บ UserID ไว้เหมือนเดิม ไม่งั้น user จะหายไปจาก income", func(t *testing.T) {
		// ส่งผลให้ที่หน้า web user เห็นว่ายังไ่ม่ได้กรอก income และจะกรอกซ้ำ
		// ซ้ำร้าย เวลา export payroll จะมี income ซ้ำกัน (ทั้ง record ที่มีและไม่มี user)
		user := testIndividualUser1
		existingIncome := Income{
			UserID:   user.ID.Hex(),
			WorkDate: "20",
		}

		result := UpdatePayroll(user, MockIncomeReq, "", &existingIncome)

		assert.Equal(t, user.ID.Hex(), result.UserID)
		assert.Equal(t, 0.05, result.WHTRate)
	})

	t.Run("calculate corporate income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := User{
			ID:   bsonutil.MustObjectIDFromHex(uidFromSession),
			Role: "corporate",
			// ปรกติเวลารายได้เกิน 1.8 ล้าน/ปี ต้องจด VAT
			// ref: https://www.rd.go.th/fileadmin/user_upload/SMEs/infographic/SME_lv1_3.pdf
			Vat:         "Y",
			DailyIncome: "5",
		}
		req := IncomeReq{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
		}
		p := NewPayroll(uidFromSession)

		err := p.parseRequest(req, user)

		assert.NoError(t, err)
		assert.Equal(t, 5*20.0, p.dailyIncome())
		assert.Equal(t, 5*20.0*0.05, p.WitholdingTax(p.dailyIncome()))
		assert.Equal(t, 7.000000000000001, p.VAT(p.dailyIncome()))
		assert.Equal(t, 100.0+7-5, p.Net(p.dailyIncome()))
		assert.Equal(t, 10*100.0, p.specialIncome())
		assert.Equal(t, 10*100.0*0.05, p.WitholdingTax(p.specialIncome()))
		assert.Equal(t, 10*100.0*0.07, p.VAT(p.specialIncome()))
		assert.Equal(t, 1000.0+70-50, p.Net(p.specialIncome()))
	})

	t.Run("เวลา export ใช้ WHT rate ที่เซฟไว้ใน income ไม่ใช้เรทปัจจุบัน", func(t *testing.T) {
		record := Income{
			WorkDate:      "20",
			SpecialIncome: "100",
			WorkingHours:  "10",
			DailyRate:     5,
			WHTRate:       0.03,
		}

		p := NewPayrollFromIncome(record)

		assert.Equal(t, 1100.0*0.03, p.totalWHT())
	})

	t.Run("income เก่าที่ไม่มี WHTRate ตอน export ถือว่าเป็น 3%", func(t *testing.T) {
		record := Income{
			WorkDate:  "20",
			DailyRate: 5,
		}

		p := NewPayrollFromIncome(record)

		assert.Equal(t, 100.0*0.03, p.totalWHT())
	})
}
