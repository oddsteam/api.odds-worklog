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
	t.Run("income created event can has more fields which worklog ignores", func(t *testing.T) {
		incomeCreatedEvent := friendslog.CreateEvent(1, "Chi", "Sweethome", 750, 20,
			"123456789122", "+66912345678", "987654321",
			15375.0, 14913.75, 750.0, 461.25, "2024-07-26T06:26:25.531Z",
			"ba1357eb-20aa-4897-9759-658bf75e8429", "user1@example.com")

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
