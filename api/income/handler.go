package income

import (
	"errors"
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
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
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /incomes [post]
func (h *HttpHandler) AddIncome(c echo.Context) error {
	var income models.IncomeReq
	if err := c.Bind(&income); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}
	if ok, err := isRequestValid(&income); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	res, err := h.Usecase.AddIncome(&income, claims.User)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
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
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /incomes/{id} [put]
func (h *HttpHandler) UpdateIncome(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	var req models.IncomeReq
	if err := c.Bind(&req); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&req); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)

	res, err := h.Usecase.UpdateIncome(id, &req, claims.User)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// GetCorporateIncomeStatus godoc
// @Summary Get Corporate Income Status List
// @Description Get Income Status List
// @Tags incomes
// @Accept  json
// @Produce  json
// @Success 200 {array} models.IncomeStatus
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /incomes/status/corporate [get]
func (h *HttpHandler) GetCorporateIncomeStatus(c echo.Context) error {
	status, err := h.Usecase.GetIncomeStatusList("Y")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, status)
}

// GetIndividualIncomeStatus godoc
// @Summary Get Individual Income Status List
// @Description Get Individual Income Status List
// @Tags incomes
// @Accept  json
// @Produce  json
// @Success 200 {array} models.IncomeStatus
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /incomes/status/individual [get]
func (h *HttpHandler) GetIndividualIncomeStatus(c echo.Context) error {
	status, err := h.Usecase.GetIncomeStatusList("N")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, status)
}

// GetIncomeByUserIdAndCurrentMonth godoc
// @Summary Get Income Of Current Month By User Id
// @Description Get Income Of Current Month By User Id
// @Tags incomes
// @Accept  json
// @Produce  json
// @Param  id path string true "User ID"
// @Success 200 {object} models.Income
// @Failure 400 {object} utils.HTTPError
// @Router /incomes/month/{id} [get]
func (h *HttpHandler) GetIncomeByUserIdAndCurrentMonth(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	income, _ := h.Usecase.GetIncomeByUserIdAndCurrentMonth(id)
	if income == nil {
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusOK, income)
}

// GetExportPdf godoc
// @Summary Get Export Pdf
// @Description Get Export to Pdf file.
// @Tages incomes
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} utils.HTTPError
// @Router /incomes/export/pdf [get]
func (h *HttpHandler) GetExportPdf(c echo.Context) error {
	filename, err := h.Usecase.ExportPdf()

	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// GetExportCorporate godoc
// @Summary Get Corporate Export Income
// @Description Get Corporate Export Income to csv file.
// @Tags incomes
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Failure 500 {object} utils.HTTPError
// @Router /incomes/export/corporate [get]
func (h *HttpHandler) GetExportCorporate(c echo.Context) error {
	checkUser, message := IsUserAdmin(c)
	if !checkUser {
		return c.JSON(http.StatusUnauthorized, message)
	}
	filename, err := h.Usecase.ExportIncome("Y")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// GetExportIndividual godoc
// @Summary Get Individual Export Income
// @Description Get Individual Export Income to csv file.
// @Tags incomes
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Failure 500 {object} utils.HTTPError
// @Router /incomes/export/individual [get]
func (h *HttpHandler) GetExportIndividual(c echo.Context) error {
	checkUser, message := IsUserAdmin(c)
	if !checkUser {
		return c.JSON(http.StatusUnauthorized, message)
	}
	filename, err := h.Usecase.ExportIncome("N")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

func (h *HttpHandler) DropIncome(c echo.Context) error {
	if err := h.Usecase.DropIncome(); err != nil {
		return utils.NewError(c, http.StatusOK, errors.New("DropIncome Failed!"))
	}
	return c.JSON(http.StatusOK, models.CommonResponse{Message: "DropIncome Success!"})
}

func IsUserAdmin(c echo.Context) (bool, string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	if claims.User.Email == "jin@odds.team" {
		return true, ""
	}
	return false, "ไม่มีสิทธิในการใช้งาน"
}
func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	incomeRepo := NewRepository(session)
	userRepo := user.NewRepository(session)
	uc := NewUsecase(incomeRepo, userRepo)
	handler := &HttpHandler{uc}

	r = r.Group("/incomes")
	r.POST("", handler.AddIncome)
	r.PUT("/:id", handler.UpdateIncome)
	r.GET("/status/corporate", handler.GetCorporateIncomeStatus)
	r.GET("/status/individual", handler.GetIndividualIncomeStatus)
	r.GET("/month/:id", handler.GetIncomeByUserIdAndCurrentMonth)
	r.GET("/export/corporate", handler.GetExportCorporate)
	r.GET("/export/individual", handler.GetExportIndividual)
	r.GET("/export/pdf", handler.GetExportPdf)
	r.DELETE("", handler.DropIncome)
}
