package income

import (
	"testing"

	"github.com/stretchr/testify/assert"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
)

func TestModelAddIncome(t *testing.T) {
	t.Run("เวลา Add income ควร save ชื่อ นามสกุล เลขบัญชี ด้วย ตอน export จะได้ไม่ต้องคำนวนแล้ว", func(t *testing.T) {
		user := userMock.User
		uidFromSession := "5bbcf2f90fd2df527bc39539"
		i := NewIncome(uidFromSession)

		res, err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.Equal(t, uidFromSession, res.UserID)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.UserID, res.UserID)
		assert.Equal(t, "116400.00", res.NetIncome)
		assert.Equal(t, "97000.00", res.NetDailyIncome)
		assert.Equal(t, "19400.00", res.NetSpecialIncome)
		assert.Equal(t, "", res.VAT)
		assert.Equal(t, "3600.00", res.WHT)
		assert.Equal(t, "120000.00", res.TotalIncome)
	})
}
