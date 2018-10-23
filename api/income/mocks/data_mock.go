package mocks

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

var (
	MockIncome = models.Income{
		TotalIncome: "123456789",
		Reason:      "เงินเดือนเดือนนี้",
		CreateBy:    "jin@odds.team",
	}

	addIncomeBy, _ = json.Marshal(MockIncome)
	AddIncomeJson  = string(addIncomeBy)
)
