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

func UpdateUser(c echo.Context) error {
	var session *mgo.Session
	var err error
	session, err = mgo.Dial("mongodb:27017")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		return err
	}
	// User := User{
	// 	FullName:          c.FormValue("fullname"),
	// 	Email:             c.FormValue("email"),
	// 	BankAccountName:   c.FormValue("bankAccountName"),
	// 	BankAccountNumber: c.FormValue("bankAccountNumber"),
	// 	TotalIncome:       c.FormValue("totalIncome"),
	// 	SubmitDate:        c.FormValue("submitDate"),
	// 	CardNumber:        c.FormValue("cardNumber"),
	// }
	colQuerier := bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}
	change := bson.M{"$set": bson.M{"fullname": "sdasdasd", "email": c.FormValue("email"), "bankAccountName": c.FormValue("bankAccountName"), "bankAccountNumber": c.FormValue("bankAccountNumber"),
		"totalIncome": c.FormValue("totalIncome"), "submitDate": c.FormValue("submitDate"), "cardNumber": c.FormValue("cardNumber")}}
	// change := bson.M{"$set": &User}
	err = session.DB("worklog-odds").C("user").Update(colQuerier, change)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, err)
}
