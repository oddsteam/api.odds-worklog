package mocks

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	MockUser = models.User{
		ID:                bson.ObjectIdHex("5bbcf2f90fd2df527bc39539"),
		FullName:          "นายทดสอบชอบลงทุน",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		CorporateFlag:     "Y",
	}

	MockToken = models.Token{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY",
	}

	MockUserById = models.User{
		ID:                "1234567890",
		FullName:          "นายทดสอบชอบลงทุน",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		CorporateFlag:     "Y",
	}

	userByte, _ = json.Marshal(MockUser)
	UserJson    = string(userByte)
	LoginJson   = `{"id": "5bbcf2f90fd2df527bc39539"}`
	Login       = models.Login{ID: "5bbcf2f90fd2df527bc39539"}

	MockUsers = []*models.User{
		{
			ID:                "1234567890",
			FullName:          "นายทดสอบชอบลงทุน",
			Email:             "test@abc.com",
			BankAccountName:   "ทดสอบชอบลงทุน",
			BankAccountNumber: "123123123123",
			ThaiCitizenID:     "1234567890123",
			CorporateFlag:     "Y",
		},
		{
			ID:                "1234567890",
			FullName:          "นายไม่ชอบลงทุน",
			Email:             "test@abc.com",
			BankAccountName:   "ทดสอบชอบลงทุน",
			BankAccountNumber: "123123123123",
			ThaiCitizenID:     "1234567890123",
			CorporateFlag:     "Y",
		},
	}
)
