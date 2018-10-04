package service

import (
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

func DeleteUser(c echo.Context) error {

	var session *mgo.Session
	var err error
	session, err = mgo.Dial(Config.DB.Host)
	defer session.Close()
	if err != nil {
		return err
	}
	err = session.DB("worklog-odds").C("user").Remove(bson.M{"_id": bson.ObjectIdHex(c.Param("id"))})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Delete Success")
}
