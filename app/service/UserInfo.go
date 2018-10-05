package service

import (
	"net/http"

	"github.com/labstack/echo"
)

func UserInfo(c echo.Context) error {
	User := []User{
		{FullName: "นายทดสอบชอบลงทุน",
			Email:             "test@abc.com",
			BankAccountName:   "ทดสอบชอบลงทุน",
			BankAccountNumber: "123123123123",
			TotalIncome:       "123123123",
			SubmitDate:        "12/12/2561",
			CardNumber:        "123123123123"},
		{FullName: "นายทดสอบชอบไม่ลงทุน",
			Email:             "test123123@abc.com",
			BankAccountName:   "ทดสอบชอบลงทุน",
			BankAccountNumber: "123123123123",
			TotalIncome:       "123123123",
			SubmitDate:        "12/12/2561",
			CardNumber:        "123123123123"},
	}
	return c.JSON(http.StatusOK, User)
}
