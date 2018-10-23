package mocks

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

var (
	MockUser = models.User{
		FullName:          "นายทดสอบชอบลงทุน",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		CoperateFlag:      "Y",
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
		CoperateFlag:      "Y",
	}

	userByte, _ = json.Marshal(MockUser)
	UserJson    = string(userByte)
	LoginJson   = `{"username": "admin", "password":"admin"}`
	Login       = models.Login{
		Username: "root",
		Password: "1234",
	}

	MockUsers = []*models.User{
		{
			ID:                "1234567890",
			FullName:          "นายทดสอบชอบลงทุน",
			Email:             "test@abc.com",
			BankAccountName:   "ทดสอบชอบลงทุน",
			BankAccountNumber: "123123123123",
			ThaiCitizenID:     "1234567890123",
			CoperateFlag:      "Y",
		},
		{
			ID:                "1234567890",
			FullName:          "นายไม่ชอบลงทุน",
			Email:             "test@abc.com",
			BankAccountName:   "ทดสอบชอบลงทุน",
			BankAccountNumber: "123123123123",
			ThaiCitizenID:     "1234567890123",
			CoperateFlag:      "Y",
		},
	}
)
