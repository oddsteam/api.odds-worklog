package models

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
)

// Mock user data - duplicated here to avoid import cycle with api/user/mock
var (
	mockUser = User{
		ID:                bson.ObjectIdHex("5bbcf2f90fd2df527bc39539"),
		Role:              "corporate",
		FirstName:         "Tester",
		LastName:          "Super",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		Vat:               "N",
		SlackAccount:      "test@abc.com",
		DailyIncome:       "5000",
		StatusTavi:        true,
		Address:           "every Where",
		StartDate:         "2022-01-01",
	}
)

var (
	MockIncome = Income{
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
	MockIncome2 = Income{
		ID:               bson.ObjectIdHex("5bd1fda30fd2df2a3e41e570"),
		UserID:           "5bbcf2f90fd2df527bc39530",
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
	MockIndividualIncome = Income{
		ID:               bson.ObjectIdHex("5bd1fda30fd2df2a3e41e568"),
		UserID:           "5bbcf2f90fd2df527bc39531",
		TotalIncome:      "110.00",
		NetIncome:        "106.70",
		NetDailyIncome:   "97.00",
		NetSpecialIncome: "9.70",
		SubmitDate:       time.Date(2022, time.Month(12), 1, 13, 30, 0, 0, time.UTC),
		Note:             "note",
		VAT:              "",
		WHT:              "3.30",
		WorkDate:         "1",
		SpecialIncome:    "10",
		WorkingHours:     "1",
		ExportStatus:     false,
	}
	MockIncomeNoNetSpecialIncome = Income{
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
	MockSoloCorporateIncome = Income{
		ID:                bson.ObjectIdHex("5bd1fda30fd2df2a3e41e568"),
		UserID:            "5bbcf2f90fd2df527bc39531",
		Name:              "บจก. โซโล่ เลเวลลิ่ง",
		TotalIncome:       "0.00",
		NetIncome:         "52600.00",
		NetDailyIncome:    "0.00",
		NetSpecialIncome:  "0.0",
		SubmitDate:        time.Date(2022, time.Month(12), 1, 13, 30, 0, 0, time.UTC),
		Note:              "note",
		VAT:               "Y",
		WHT:               "0.00",
		WorkDate:          "20",
		DailyRate:         2630,
		SpecialIncome:     "0",
		WorkingHours:      "0",
		ExportStatus:      false,
		IsVATRegistered:   true,
		BankAccountNumber: "2462737202",
	}
	MockSwardCorporateIncome = Income{
		ID:                bson.ObjectIdHex("5bd1fda30fd2df2a3e41e568"),
		UserID:            string(bson.ObjectIdHex("5bbcf2f90fd2df527bc39531")),
		Name:              "บจ. ดาบพิฆาตอสูร",
		TotalIncome:       "0.00",
		NetIncome:         "5260.00",
		NetDailyIncome:    "0.00",
		NetSpecialIncome:  "0.0",
		SubmitDate:        time.Date(2022, time.Month(12), 1, 13, 30, 0, 0, time.UTC),
		Note:              "note",
		VAT:               "Y",
		WHT:               "0.00",
		WorkDate:          "20",
		DailyRate:         263,
		SpecialIncome:     "0",
		WorkingHours:      "0",
		ExportStatus:      false,
		IsVATRegistered:   true,
		BankAccountNumber: "1102480447",
	}
	MockIncomeReq = IncomeReq{
		WorkDate:      "20",
		Note:          "ข้อมูลที่อยากบอก",
		SpecialIncome: "2000",
		WorkingHours:  "10",
	}

	MockIncomeStatus = IncomeStatus{
		User:       &mockUser,
		SubmitDate: "2018-10-24 20:30:40",
		Status:     "Y",
		WorkDate:   "20",
	}

	MockCorporateIncomeStatus = IncomeStatus{
		User:   &mockUser,
		Status: "Y",
	}
	MockIndividualIncomeStatus = IncomeStatus{
		User:   &mockUser,
		Status: "N",
	}
	MockIncomeList                   = []*Income{&MockIncome, &MockIncome2}
	MockIncomeListNoNetSpecialIncome = []*Income{&MockIncomeNoNetSpecialIncome, &MockIncome2}
	MockIncomeStatusList             = []*IncomeStatus{&MockIncomeStatus}
	IncomeByte, _                    = json.Marshal(MockIncome)
	MockIncomeJson                   = string(IncomeByte)

	IncomeReqByte, _  = json.Marshal(MockIncomeReq)
	MockIncomeReqJson = string(IncomeReqByte)

	IncomeResByte, _  = json.Marshal(MockIncomeStatus)
	MockIncomeResJson = string(IncomeResByte)
)

