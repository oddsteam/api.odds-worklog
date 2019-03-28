package mock_income

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"

	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

var (
	MockIncome = models.Income{
		ID:            bson.ObjectIdHex("5bd1fda30fd2df2a3e41e569"),
		UserID:        "5bbcf2f90fd2df527bc39539",
		TotalIncome:   "100000",
		NetIncome:     "100000",
		SubmitDate:    time.Now(),
		Note:          "ข้อมูลที่อยากบอก",
		VAT:           "7000",
		WHT:           "3000",
		WorkDate:      "20",
		SpecialIncome: "2000",
	}
	MockIncomeReq = models.IncomeReq{
		WorkDate:      "100000",
		Note:          "ข้อมูลที่อยากบอก",
		SpecialIncome: "5000",
	}

	MockIncomeStatus = models.IncomeStatus{
		User:       &userMock.User,
		SubmitDate: "2018-10-24 20:30:40",
		Status:     "Y",
		WorkDate:   "20",
	}

	MockCorporateIncomeStatus = models.IncomeStatus{
		User:   &userMock.User,
		Status: "Y",
	}
	MockIndividualIncomeStatus = models.IncomeStatus{
		User:   &userMock.User,
		Status: "N",
	}
	MockIncomeStatusList = []*models.IncomeStatus{&MockIncomeStatus}
	IncomeByte, _        = json.Marshal(MockIncome)
	MockIncomeJson       = string(IncomeByte)

	IncomeReqByte, _  = json.Marshal(MockIncomeReq)
	MockIncomeReqJson = string(IncomeReqByte)

	IncomeResByte, _  = json.Marshal(MockIncomeStatus)
	MockIncomeResJson = string(IncomeResByte)
)
