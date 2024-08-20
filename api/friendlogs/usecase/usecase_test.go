package usecase_test

import (
	"fmt"
	"testing"
	"time"

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
		createdAt := time.Date(2024, time.Month(7), 26, 6, 26, 25, 531000000, time.UTC)
		assert.Equal(t, createdAt, income.SubmitDate)
		assert.Equal(t, createdAt, income.LastUpdate)
	})
	t.Run("The total amount of the Income which a Coop added in friendslog is calculated", func(t *testing.T) {
		workDate := 20
		incomeCreatedEvent := simpleCoopIncomeEvent(workDate, 750)

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "15000.00", income.TotalIncome)
		assert.Equal(t, "450.00", income.WHT)
		assert.Equal(t, "14550.00", income.NetIncome)
	})

	t.Run("income contains note when it was added", func(t *testing.T) {
		incomeCreatedEvent := addedIncomeEventAt(20, "2024-07-26T06:26:25.531Z")

		income := u.AddIncome(incomeCreatedEvent)

		assert.Equal(t, "Added on 2024-07-26T06:26:25.531Z", income.Note)
	})
	t.Run("worklog ignores irrelevant fields from friendslog", func(t *testing.T) {
		incomeCreatedEvent := `{
			"otherField": "value",
			"income":{
				"workDate":20
			},
			"registration":{
				"thai_citizen_id":"0123456789121",
				"daily_income":"750",
				"userId":"userId"
			}
		}`

		income := u.AddIncome(incomeCreatedEvent)

		assert.NotNil(t, income)
	})
}

func TestUsecaseUpdateIncome(t *testing.T) {
	u := usecase.NewUsecase()
	t.Run("Coop updated income from 20 -> 21 in friendslog success", func(t *testing.T) {
		allAddedIncomes := givenThereIsAnIncomeExist("20", 750, "Added on 2024-07-22T06:26:25.531Z")
		incomeUpdatedEvent := updatedIncomeEventAt(21, "2024-07-23T06:26:25.531Z")

		income := u.UpdateIncome(allAddedIncomes, incomeUpdatedEvent)

		expectedNote := "Added on 2024-07-22T06:26:25.531Z\nUpdated on 2024-07-23T06:26:25.531Z"
		assert.Equal(t, "21", income.WorkDate)
		assert.Equal(t, expectedNote, income.Note)
	})

	t.Run("only update income of the same citizen id", func(t *testing.T) {
		allAddedIncomes := []*models.Income{
			{
				ID:            "1",
				ThaiCitizenID: "another id", WorkDate: "20", DailyRate: 750,
			},
			{
				ID:            "2",
				ThaiCitizenID: "0123456789121",
				WorkDate:      "20",
				DailyRate:     750,
				Note:          "Added on 2024-07-22T06:26:25.531Z",
			},
		}
		incomeUpdatedEvent := updatedIncomeEventAt(21, "2024-07-23T06:26:25.531Z")

		income := u.UpdateIncome(allAddedIncomes, incomeUpdatedEvent)

		assert.Equal(t, bson.ObjectId("2"), income.ID)
	})
}

func givenThereIsAnIncomeExist(days string, rate float64, n string) []*models.Income {
	return []*models.Income{
		{
			ThaiCitizenID: "0123456789121",
			WorkDate:      days,
			DailyRate:     rate,
			Note:          n,
		},
	}
}

func addedIncomeEventAt(days int, createdAt string) string {
	return fullCoopIncomeEvent("Chi", "Sweethome", 750, days,
		"0123456789121", "+66912345678", "987654321",
		createdAt, createdAt, "user1@example.com")
}
func updatedIncomeEventAt(days int, updatedAt string) string {
	return fullCoopIncomeEvent("Chi", "Sweethome", 750, days,
		"0123456789121", "+66912345678", "987654321",
		"", updatedAt, "user1@example.com")
}

func fullCoopIncomeEvent(firstName string, lastName string,
	dailyRate float64, workDays int, thaiCitizenID string,
	phone string, bankAcocuntNumber string, createAt string, updatedAt string, email string) string {

	return usecase.CreateEvent(1, firstName, lastName, int(dailyRate), workDays,
		thaiCitizenID, phone, bankAcocuntNumber,
		0, 0, 0, 0, createAt, updatedAt, "friendslogId", email)
}

func simpleCoopIncomeEvent(workDate int, dailyRate float64) string {
	return fmt.Sprintf(`{
			"income":{
				"workDate":%d
			},
			"registration":{
				"thai_citizen_id":"0123456789121",
				"daily_income":"%f",
				"userId":"userId"
			}
		}`, workDate, dailyRate)
}
