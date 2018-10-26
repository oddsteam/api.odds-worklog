package login

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	r.POST("/login", func(c echo.Context) error {
		return login(c, session)
	})
}

func login(c echo.Context, session *mongo.Session) error {
	var u models.Login
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	user, err := getUserByID(u.ID, session)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	user.BankAccountName = ""
	user.BankAccountNumber = ""
	user.ThaiCitizenID = ""

	claims := &models.JwtCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("GmkZGF3CmpZNs88dLvbV"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	tk := &models.Token{
		Token: t,
	}
	return c.JSON(http.StatusOK, tk)
}

func getUserByID(id string, session *mongo.Session) (*models.User, error) {
	user := new(models.User)
	coll := session.GetCollection("user")
	err := coll.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
