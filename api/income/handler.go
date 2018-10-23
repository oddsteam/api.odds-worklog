package income

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpHandler struct {
	Usecase Usecase
}

func isRequestValid(m *models.Income) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

func (h *HttpHandler) AddIncome(c echo.Context) error {
	var addIncome models.Income
	if err := c.Bind(&addIncome); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&addIncome); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	addIncome.UserId = claims.Name
	income, err := h.Usecase.AddIncome(&addIncome)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, income)
}

func NewHttpHandler(e *echo.Echo, session *mongo.Session) {
	ur := newRepository(session)
	uc := newUsecase(ur)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handler := &HttpHandler{uc}
	config := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r := e.Group("/")
	r.Use(middleware.JWTWithConfig(config))
	r.POST("income", handler.AddIncome)
}
