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
		TotalIncome:       "123123123",
		SubmitDate:        "12/12/2561",
		ThaiCitizenID:     "1234567890123",
	}

	MockToken = models.Token{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY",
	}

	userByte, _ = json.Marshal(MockUser)
	UserJson    = string(userByte)
	LoginJson   = `{"username": "root", "password":"1234"}`
)
