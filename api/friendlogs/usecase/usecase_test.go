package usecase_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/friendlogs/usecase"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestUsecaseAddIncome(t *testing.T) {
	u := usecase.NewUsecase()
	t.Run("Income which a Coop added in friendslog is saved in worklog", func(t *testing.T) {
		thaiCitizenID := "0123456789121"
		incomeCreatedEvent := fullCoopIncomeEvent("Chi", "Sweethome", 750, 20,
			thaiCitizenID, "+66912345678", "987654321",
			"2024-07-26T06:26:25.531Z", "", "user1@example.com")

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "Chi Sweethome", income.Name)
		assert.Equal(t, "20", income.WorkDate)
		assert.Equal(t, thaiCitizenID, income.ThaiCitizenID)
		assert.Equal(t, "+66912345678", income.Phone)
		assert.Equal(t, "987654321", income.BankAccountNumber)
		assert.Equal(t, 750.0, income.DailyRate)
		assert.Equal(t, "Chi Sweethome", income.BankAccountName)
		assert.Equal(t, "user1@example.com", income.Email)
	})
	t.Run("The total amount of the Income which a Coop added in friendslog is calculated", func(t *testing.T) {
		workDate := 20
		incomeCreatedEvent := simpleCoopIncomeEvent("0123456789121", workDate, 750)

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "15000.00", income.TotalIncome)
		assert.Equal(t, "450.00", income.WHT)
		assert.Equal(t, "14550.00", income.NetIncome)
	})

	t.Run("income contains note when it was added", func(t *testing.T) {
		incomeCreatedEvent := fullCoopIncomeEvent("Chi", "Sweethome", 750, 20,
			"0123456789121", "+66912345678", "987654321",
			"2024-07-26T06:26:25.531Z", "", "user1@example.com")

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "Added on 2024-07-26T06:26:25.531Z", income.Note)
	})
	t.Run("income created event can has more fields which worklog ignores", func(t *testing.T) {
		incomeCreatedEvent := usecase.CreateEvent(1, "Chi", "Sweethome", 750, 20,
			"123456789122", "+66912345678", "987654321",
			15375.0, 14913.75, 750.0, 461.25, "2024-07-26T06:26:25.531Z", "",
			"ba1357eb-20aa-4897-9759-658bf75e8429", "user1@example.com")

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "15000.00", income.TotalIncome)
	})
}

func TestUsecaseUpdateIncome(t *testing.T) {
	u := usecase.NewUsecase()
	t.Run("Coop updated income from 20 -> 21 in friendslog success", func(t *testing.T) {
		thaiCitizenID := "0123456789121"
		allAddedIncomes := []*models.Income{
			{
				ThaiCitizenID: thaiCitizenID,
				WorkDate:      "20",
				DailyRate:     750,
				Note:          "Added on 2024-07-22T06:26:25.531Z",
			},
		}

		incomeUpdatedEvent := fullCoopIncomeEvent("Chi", "Sweethome", 750, 21,
			thaiCitizenID, "+66912345678", "987654321",
			"", "2024-07-23T06:26:25.531Z", "user1@example.com")

		income := u.UpdateIncome(allAddedIncomes, incomeUpdatedEvent)

		expectedNote := []string{
			"Added on 2024-07-22T06:26:25.531Z",
			"Updated on 2024-07-23T06:26:25.531Z",
		}
		assert.Equal(t, "21", income.WorkDate)
		assert.Equal(t, strings.Join(expectedNote, "\n"), income.Note)
	})

	t.Run("only update income of the same citizen id", func(t *testing.T) {
		thaiCitizenID := "0123456789121"
		allAddedIncomes := []*models.Income{
			{
				ID:            "1",
				ThaiCitizenID: "another id", WorkDate: "20", DailyRate: 750,
			},
			{
				ID:            "2",
				ThaiCitizenID: thaiCitizenID,
				WorkDate:      "20",
				DailyRate:     750,
				Note:          "Added on 2024-07-22T06:26:25.531Z",
			},
		}

		incomeUpdatedEvent := fullCoopIncomeEvent("Chi", "Sweethome", 750, 21,
			thaiCitizenID, "+66912345678", "987654321",
			"", "2024-07-23T06:26:25.531Z", "user1@example.com")

		income := u.UpdateIncome(allAddedIncomes, incomeUpdatedEvent)

		assert.Equal(t, bson.ObjectId("2"), income.ID)
	})
}
func fullCoopIncomeEvent(firstName string, lastName string,
	dailyRate float64, workDays int, thaiCitizenID string,
	phone string, bankAcocuntNumber string, createAt string, updatedAt string, email string) string {

	return usecase.CreateEvent(1, firstName, lastName, int(dailyRate), workDays,
		thaiCitizenID, phone, bankAcocuntNumber,
		0, 0, 0, 0, createAt, updatedAt, "friendslogId", email)
}

func simpleCoopIncomeEvent(thaiCitizenID string, workDate int, dailyRate float64) string {
	return fmt.Sprintf(`{
			"income":{
				"workDate":%d
			},
			"registration":{
				"thai_citizen_id":"%s",
				"daily_income":"%f",
				"userId":"userId"
			}
		}`, workDate, thaiCitizenID, dailyRate)
}
