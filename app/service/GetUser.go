package service

import (
	"net/http"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
)

func GetUser(c echo.Context) error {
	user := []User{}
	var session *mgo.Session
	var err error
	session, err = mgo.Dial("mongodb:27017")
	defer session.Close()
	err = session.DB("worklog-odds").C("user").Find(bson.M{}).All(&user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)

}
