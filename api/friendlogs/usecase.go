package friendslog

import (
	"encoding/json"

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
	uid := "000000000000000000000000"
	return *income.CreateIncome(uid,
		data.Registration.ThaiCitizenID,
		data.Registration.DailyIncome, data.Income.WorkDate, "0", "0")
}

type IncomeCreatedEvent struct {
	Income       Income       `json:"income"`
	Registration Registration `json:"registration"`
}

type Income struct {
	WorkDate string `json:"workDate"`
}

type Registration struct {
	ThaiCitizenID string `json:"thai_citizen_id"`
	DailyIncome   string `json:"daily_income"`
	UserID        string `json:"userId"`
}
