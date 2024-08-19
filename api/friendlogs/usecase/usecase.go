package usecase

import (
	"encoding/json"
	"fmt"

	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct{}

func NewUsecase() *usecase {
	return &usecase{}
}

func (u *usecase) AddIncome(incomeStr string) models.Income {
	var data IncomeCreatedEvent
	err := json.Unmarshal([]byte(incomeStr), &data)
	utils.FailOnError(err, "Error parsing JSON")
	user := dataToUser(data)
	req := models.IncomeReq{
		WorkDate:      fmt.Sprint(data.Income.WorkDate),
		SpecialIncome: "0",
		WorkingHours:  "0",
	}
	note := fmt.Sprintf("Added on %s", data.Income.CreatedAt)
	return *income.CreateIncome(user, req, note)
}

func (u *usecase) UpdateIncome(original models.Income, incomeStr string) models.Income {
	var data IncomeCreatedEvent
	err := json.Unmarshal([]byte(incomeStr), &data)
	utils.FailOnError(err, "Error parsing JSON")
	user := dataToUser(data)
	req := models.IncomeReq{
		WorkDate:      fmt.Sprint(data.Income.WorkDate),
		SpecialIncome: "0",
		WorkingHours:  "0",
	}
	note := fmt.Sprintf("%s\nUpdated on %s", original.Note, data.Income.UpdatedAt)
	return *income.CreateIncome(user, req, note)
}

func dataToUser(data IncomeCreatedEvent) models.User {
	uid := "000000000000000000000000"
	user := income.GivenIndividualUser(uid, data.Registration.DailyIncome)
	user.ThaiCitizenID = data.Registration.ThaiCitizenID
	user.FirstName = data.Registration.FirstName
	user.LastName = data.Registration.LastName
	user.Phone = data.Registration.Phone
	user.BankAccountNumber = data.Registration.BankAccountNumber
	user.BankAccountName = user.GetFullname()
	user.Email = data.Registration.Email
	return user
}

type IncomeCreatedEvent struct {
	Income       Income       `json:"income"`
	Registration Registration `json:"registration"`
}

type Income struct {
	WorkDate  int    `json:"workDate"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Registration struct {
	ThaiCitizenID     string `json:"thai_citizen_id"`
	DailyIncome       string `json:"daily_income"`
	UserID            string `json:"userId"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Phone             string `json:"phone"`
	BankAccountNumber string `json:"bank_no"`
	Email             string `json:"email"`
}
