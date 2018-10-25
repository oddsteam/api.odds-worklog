package income

import (
	"errors"
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"

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

func isRequestValid(m *models.IncomeReq) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

func (h *HttpHandler) AddIncome(c echo.Context) error {
	var income models.IncomeReq
	if err := c.Bind(&income); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&income); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	res, err := h.Usecase.AddIncome(&income, claims.User)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *HttpHandler) UpdateIncome(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errors.New("invalid path"))
	}

	var req models.IncomeReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&req); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)

	res, err := h.Usecase.UpdateIncome(id, &req, claims.User)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func NewHttpHandler(e *echo.Echo, config middleware.JWTConfig, session *mongo.Session) {
	incomeRepo := newRepository(session)
	userRepo := user.NewRepository(session)
	uc := newUsecase(incomeRepo, userRepo)
	handler := &HttpHandler{uc}

	r := e.Group("/incomes")
	r.Use(middleware.JWTWithConfig(config))

	r.POST("", handler.AddIncome)
	r.PUT("/:id", handler.UpdateIncome)
}
