package income

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestModelAddIncome(t *testing.T) {
	t.Run("เวลา Add income ควร save ชื่อบัญชี เลขบัญชี และจำนวนเงินด้วย ตอน export จะได้ไม่ต้องคำนวนแล้ว", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.Equal(t, uidFromSession, res.UserID)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.UserID, res.UserID)
		assert.Equal(t, "58200.00", res.NetIncome)
		assert.Equal(t, "38800.00", res.NetDailyIncome)
		assert.Equal(t, "19400.00", res.NetSpecialIncome)
		assert.Equal(t, "", res.VAT)
		assert.Equal(t, "1800.00", res.WHT)
		assert.Equal(t, "60000.00", res.TotalIncome)
		assert.Equal(t, user.BankAccountName, res.BankAccountName)
		assert.Equal(t, user.BankAccountNumber, res.BankAccountNumber)
	})

	t.Run("เวลา Add income ควร save ชื่อ นามสกุล เลขบัตรประชาชนเวลา export ให้บัญชี เค้าจะได้รู้ว่าจ่ายเงินให้ใคร", func(t *testing.T) {
		user := userMock.IndividualUser1
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.ThaiCitizenID, res.ThaiCitizenID)
		assert.Equal(t, "first last", res.Name)
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

	t.Run("calculate individual income", func(t *testing.T) {
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		user := models.User{
			ID:          bson.ObjectIdHex(uidFromSession),
			Role:        "individual",
			Vat:         "N",
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
		assert.Equal(t, 0.0, i.VAT(i.totalIncome()))
		assert.Equal(t, 100.0+0-3, i.Net(i.totalIncome()))
		assert.Equal(t, 10*100.0, i.specialIncome())
		assert.Equal(t, 10*100.0*0.03, i.WitholdingTax(i.specialIncome()))
		assert.Equal(t, 0.0, i.VAT(i.specialIncome()))
		assert.Equal(t, 1000.0+0-30, i.Net(i.specialIncome()))
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
