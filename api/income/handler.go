package income

import (
	"errors"
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/httputil"
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

// AddIncome godoc
// @Summary Add Income
// @Description Add Income
// @Tags incomes
// @Accept  json
// @Produce  json
// @Param income body models.IncomeReq true  " "
// @Success 200 {object} models.Income
// @Failure 400 {object} httputil.HTTPError
// @Failure 422 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /incomes [post]
func (h *HttpHandler) AddIncome(c echo.Context) error {
	var income models.IncomeReq
	if err := c.Bind(&income); err != nil {
		return httputil.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&income); !ok {
		return httputil.NewError(c, http.StatusBadRequest, err)
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	res, err := h.Usecase.AddIncome(&income, claims.User)
	if err != nil {
		return httputil.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateIncome godoc
// @Summary Update Income
// @Description Update Income
// @Tags incomes
// @Accept  json
// @Produce  json
// @Param income body models.IncomeReq true  " "
// @Param  id path string true "Income ID"
// @Success 200 {object} models.Income
// @Failure 400 {object} httputil.HTTPError
// @Failure 422 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /incomes/{id} [put]
func (h *HttpHandler) UpdateIncome(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return httputil.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	var req models.IncomeReq
	if err := c.Bind(&req); err != nil {
		return httputil.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&req); !ok {
		return httputil.NewError(c, http.StatusBadRequest, err)
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)

	res, err := h.Usecase.UpdateIncome(id, &req, claims.User)
	if err != nil {
		return httputil.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// GetCorporateIncomeStatus godoc
// @Summary Get Corporate Income Status
// @Description Get Income Status
// @Tags incomes
// @Accept  json
// @Produce  json
// @Success 200 {array} models.IncomeStatus
// @Failure 400 {object} httputil.HTTPError
// @Failure 422 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /incomes/status/corporate [get]
func (h *HttpHandler) GetCorporateIncomeStatus(c echo.Context) error {
	users, err := h.Usecase.GetIncomeStatusList("Y")
	if err != nil {
		return httputil.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

// GetIncomeByUserIdAndCurrentMonth godoc
// @Summary Get Income Of Current Month By User Id
// @Description Get Income Of Current Month By User Id
// @Tags incomes
// @Accept  json
// @Produce  json
// @Param  id path string true "User ID"
// @Success 200 {object} models.Income
// @Failure 400 {object} httputil.HTTPError
// @Failure 422 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /incomes/month/{id} [get]
func (h *HttpHandler) GetIncomeByUserIdAndCurrentMonth(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return httputil.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	income, err := h.Usecase.GetIncomeByUserIdAndCurrentMonth(id)
	if income == nil {
		return c.JSON(http.StatusOK, nil)
	}
	if err != nil {
		return httputil.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, income)
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	incomeRepo := newRepository(session)
	userRepo := user.NewRepository(session)
	uc := newUsecase(incomeRepo, userRepo)
	handler := &HttpHandler{uc}

	r = r.Group("/incomes")
	r.POST("", handler.AddIncome)
	r.PUT("/:id", handler.UpdateIncome)
	r.GET("/status/corporate", handler.GetCorporateIncomeStatus)
	r.GET("/month/:id", handler.GetIncomeByUserIdAndCurrentMonth)
}
