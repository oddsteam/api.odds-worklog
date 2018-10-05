package service

import (
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

//----------
// Handlers
//----------

func InsertUser(c echo.Context) error {
	var session *mgo.Session
	var err error
	session, err = mgo.Dial("mongodb:27017")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		return err
	}

	newId := bson.NewObjectId()

	user := User{
		ID:                newId,
		FullName:          c.FormValue("fullname"),
		Email:             c.FormValue("email"),
		BankAccountName:   c.FormValue("bankAccountName"),
		BankAccountNumber: c.FormValue("bankAccountNumber"),
		TotalIncome:       c.FormValue("totalIncome"),
		SubmitDate:        c.FormValue("submitDate"),
		CardNumber:        c.FormValue("cardNumber"),
	}

	err = session.DB("worklog-odds").C("user").Insert(user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}
