package mocks

import (
	"encoding/json"

	userMocks "gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

var (
	MockIncome = models.Income{
		UserID:      "5bbcf2f90fd2df527bc39539",
		TotalIncome: "100000",
		SubmitDate:  "2018-10-24 20:30:40",
		Note:        "ข้อมูลที่อยากบอก",
		VAT:         "7000",
		WHT:         "3000",
	}

	MockIncomeReq = models.IncomeReq{
		TotalIncome: "100000",
		Note:        "ข้อมูลที่อยากบอก",
	}

	MockIncomeRes = models.IncomeRes{
		User:       userMocks.MockUser,
		SubmitDate: "2018-10-24 20:30:40",
		Status:     "Y",
	}

	IncomeByte, _  = json.Marshal(MockIncome)
	MockIncomeJson = string(IncomeByte)

	IncomeReqByte, _  = json.Marshal(MockIncomeReq)
	MockIncomeReqJson = string(IncomeReqByte)

	IncomeResByte, _  = json.Marshal(MockIncomeRes)
	MockIncomeResJson = string(IncomeResByte)
)
