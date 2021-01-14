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
		ID:               bson.ObjectIdHex("5bd1fda30fd2df2a3e41e569"),
		UserID:           "5bbcf2f90fd2df527bc39539",
		TotalIncome:      "100000",
		NetIncome:        "116400.00",
		NetDailyIncome:   "97000.00",
		NetSpecialIncome: "19400.00",
		SubmitDate:       time.Now(),
		Note:             "ข้อมูลที่อยากบอก",
		VAT:              "0",
		WHT:              "3600.00",
		WorkDate:         "20",
		SpecialIncome:    "2000",
		WorkingHours:     "10",
		ExportStatus:     false,
	}
	MockIncome2 = models.Income{
		ID:               bson.ObjectIdHex("5bd1fda30fd2df2a3e41e569"),
		UserID:           "5bbcf2f90fd2df527bc39539",
		TotalIncome:      "100000",
		NetIncome:        "00.00",
		NetDailyIncome:   "48500.00",
		NetSpecialIncome: "1940.00",
		SubmitDate:       time.Now(),
		Note:             "ข้อมูลที่อยากบอก",
		VAT:              "0",
		WHT:              "1560.00",
		WorkDate:         "20",
		SpecialIncome:    "200",
		WorkingHours:     "10",
		ExportStatus:     false,
	}
	MockIncomeNoNetSpecialIncome = models.Income{
		ID:             bson.ObjectIdHex("5bd1fda30fd2df2a3e41e569"),
		UserID:         "5bbcf2f90fd2df527bc39539",
		TotalIncome:    "100000",
		NetIncome:      "00.00",
		NetDailyIncome: "48500.00",
		SubmitDate:     time.Now(),
		Note:           "ข้อมูลที่อยากบอก",
		VAT:            "0",
		WHT:            "1560.00",
		WorkDate:       "20",
		SpecialIncome:  "200",
		WorkingHours:   "10",
		ExportStatus:   false,
	}
	MockIncomeReq = models.IncomeReq{
		WorkDate:      "20",
		Note:          "ข้อมูลที่อยากบอก",
		SpecialIncome: "2000",
		WorkingHours:  "10",
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
	MockIncomeList                   = []*models.Income{&MockIncome, &MockIncome2}
	MockIncomeListNoNetSpecialIncome = []*models.Income{&MockIncomeNoNetSpecialIncome, &MockIncome2}
	MockIncomeStatusList             = []*models.IncomeStatus{&MockIncomeStatus}
	IncomeByte, _                    = json.Marshal(MockIncome)
	MockIncomeJson                   = string(IncomeByte)

	IncomeReqByte, _  = json.Marshal(MockIncomeReq)
	MockIncomeReqJson = string(IncomeReqByte)

	IncomeResByte, _  = json.Marshal(MockIncomeStatus)
	MockIncomeResJson = string(IncomeResByte)
)
