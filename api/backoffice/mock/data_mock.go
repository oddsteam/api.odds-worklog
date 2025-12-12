package mock_backoffice

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"

	siteMock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

var (
	MockUserIncome = models.UserIncome{
		ID:                bson.ObjectIdHex("5bd1fda30fd2df2a3e41e569"),
		Role:              "individual",
		FirstName:         "Apinrat",
		LastName:          "Jaidee",
		CorporateName:     "cccc",
		Email:             "apinrat@odds.team",
		BankAccountName:   "อพินรต ใจดี",
		BankAccountNumber: "123456789",
		ThaiCitizenID:     "1309901271351",
		Vat:               "N",
		SlackAccount:      "apinrat@odds.team",
		Transcript:        "",
		SiteID:            "5bbcf2f90fd2df527bc39530",
		Project:           "MMS",
		ImageProfile:      "files/images/prayuth_janogkachart_iqM8fu4U8JWX..png",
		DegreeCertificate: "",
		IDCard:            "",
		Site:              &siteMock.MockSite,
		Create:            time.Now(),
		LastUpdate:        time.Now(),
		DailyIncome:       "3000",
		Address:           "265/28",
		StatusTavi:        true,
		Incomes:           *&models.MockIncomeList,
	}

	MockBackOfficeKey = models.BackOfficeKey{
		ID:  bson.ObjectIdHex("5bd1fda30fd2df2a3e41e522"),
		Key: "TESTKEY",
	}

	MockBackOfficeKeyReq = models.BackOfficeKey{
		ID:  bson.ObjectIdHex("5bd1fda30fd2df2a3e41e533"),
		Key: "TESTKEY",
	}

	MockInvalideBackOfficeKeyReq = models.BackOfficeKey{
		ID:  bson.ObjectIdHex("5bd1fda30fd2df2a3e41e533"),
		Key: "TESTKEY1",
	}

	MockUserIncomeList = []*models.UserIncome{&MockUserIncome, &MockUserIncome}

	backOfficeKeyReqByte, _ = json.Marshal(MockBackOfficeKeyReq)
	BackOfficeKeyReqJson    = string(backOfficeKeyReqByte)

	invalidBackOfficeKeyReqByte, _ = json.Marshal(MockInvalideBackOfficeKeyReq)
	InvalidBackOfficeKeyReqJson    = string(invalidBackOfficeKeyReqByte)
)
