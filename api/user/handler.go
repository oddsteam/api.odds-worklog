package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpHandler struct {
	Usecase Usecase
}

func isRequestValid(m *models.User) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

func (h *HttpHandler) CreateUser(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&u); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.Usecase.CreateUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *HttpHandler) GetUser(c echo.Context) error {
	users, err := h.Usecase.GetUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *HttpHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.Usecase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNoContent, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	u := models.User{
		ID: bson.ObjectIdHex(id),
	}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&u); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.Usecase.UpdateUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.Usecase.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *HttpHandler) Login(c echo.Context) error {
	var u models.Login
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	user, err := h.Usecase.GetUserByID(u.ID)
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

func (h *HttpHandler) UpdatePartialUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.Usecase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	b, _ := ioutil.ReadAll(c.Request().Body)
	err = json.Unmarshal(b, &user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	newUser, err := h.Usecase.UpdateUser(user)
	return c.JSON(http.StatusOK, newUser)
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	ur := NewRepository(session)
	uc := newUsecase(ur)
	handler := &HttpHandler{uc}

	r = r.Group("/users")
	r.GET("", handler.GetUser)
	r.POST("", handler.CreateUser)
	r.GET("/:id", handler.GetUserByID)
	r.PUT("/:id", handler.UpdateUser)
	r.DELETE("/:id", handler.DeleteUser)
	r.PATCH("/:id", handler.UpdatePartialUser)
}
