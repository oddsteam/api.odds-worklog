package friendslog_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	friendslog "gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs"
)

func TestUsecaseAddIncome(t *testing.T) {
	u := friendslog.NewUsecase()
	t.Run("Income which a Coop added in friendslog is saved in worklog", func(t *testing.T) {
		thaiCitizenID := "0123456789121"
		workDate := "20"
		dailyRate := 750.0
		incomeCreatedEvent := simpleCoopIncomeEvent(thaiCitizenID, workDate, dailyRate)

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, workDate, income.WorkDate)
		assert.Equal(t, thaiCitizenID, income.ThaiCitizenID)
		assert.Equal(t, dailyRate, income.DailyRate)
	})
	t.Run("The total amount of the Income which a Coop added in friendslog is calculated", func(t *testing.T) {
		workDate := "20"
		incomeCreatedEvent := simpleCoopIncomeEvent("0123456789121", workDate, 750)

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "15000.00", income.TotalIncome)
	})
}

func simpleCoopIncomeEvent(thaiCitizenID string, workDate string, dailyRate float64) string {
	return fmt.Sprintf(`{
			"income":{
				"workDate":"%s"
			},
			"registration":{
				"thai_citizen_id":"%s",
				"daily_income":"%f",
				"userId":"userId"
			}
		}`, workDate, thaiCitizenID, dailyRate)
}
