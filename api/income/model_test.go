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

		err := i.prepareDataForAddIncome(incomeMock.MockIncomeReq, user)

		assert.NoError(t, err)
		assert.Equal(t, uidFromSession, i.UserID)
		assert.Equal(t, "116400.00", i.NetIncomeStr)
		assert.Equal(t, "97000.00", i.NetDailyIncomeStr)
		assert.Equal(t, "19400.00", i.NetSpecialIncomeStr)
		assert.Equal(t, "", i.VATStr)
		assert.Equal(t, "3600.00", i.WHTStr)
		assert.Equal(t, "120000.00", i.TotalIncomeStr)
	})
}
