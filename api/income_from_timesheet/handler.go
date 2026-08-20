package income_from_timesheet

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gitlab.odds.team/worklog/api.odds-worklog/repositories"
)

// HttpHandler serves the same export formats as /incomes, but sourced from
// income_from_timesheet instead of income.
type HttpHandler struct {
	ExportIncomeFromTimesheetUsecase usecases.ForUsingExportIncome
}

// GetExportIndividual godoc
// @Summary Export individual income_from_timesheet data to CSV
// @Description Exports income_from_timesheet records for individual users, in the same format
// @Description as the real individual income export.
// @Tags income-from-timesheet
// @Produce json
// @Param month path string true "Month index (0 = current month, else previous month)"
// @Success 200 {array} string
// @Failure 401 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /income-from-timesheet/export/individual/{month} [get]
func (h *HttpHandler) GetExportIndividual(c echo.Context) error {
	month, response, ok := h.authorizedMonthParam(c)
	if !ok {
		return response
	}
	filename, err := h.ExportIncomeFromTimesheetUsecase.ExportIncome("individual", month)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// GetExportPeakIndividual godoc
// @Summary Export individual income_from_timesheet data to a PEAK CSV
// @Tags income-from-timesheet
// @Produce json
// @Param month path string true "Month index (0 = current month, else previous month)"
// @Success 200 {array} string
// @Failure 401 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /income-from-timesheet/export/peak/individual/{month} [get]
func (h *HttpHandler) GetExportPeakIndividual(c echo.Context) error {
	month, response, ok := h.authorizedMonthParam(c)
	if !ok {
		return response
	}
	filename, err := h.ExportIncomeFromTimesheetUsecase.ExportPeak("individual", month)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// PostExport godoc
// @Summary Export individual income_from_timesheet data to CSV for a given period
// @Description Exports income_from_timesheet records between startDate and endDate (both "MM/YYYY").
// @Tags income-from-timesheet
// @Accept json
// @Produce json
// @Param body body models.ExportInComeReq true "Export period"
// @Success 200 {array} string
// @Failure 401 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /income-from-timesheet/export [post]
func (h *HttpHandler) PostExport(c echo.Context) error {
	req, response, ok := h.authorizedExportPeriod(c)
	if !ok {
		return response
	}
	filename, err := h.ExportIncomeFromTimesheetUsecase.ExportIncomeByStartDateAndEndDate(req.role, req.startDate, req.endDate)
	if err != nil {
		log.Println(err.Error())
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// PostExportPeak godoc
// @Summary Export individual income_from_timesheet data to a PEAK CSV for a given period
// @Tags income-from-timesheet
// @Accept json
// @Produce json
// @Param body body models.ExportInComeReq true "Export period"
// @Success 200 {array} string
// @Failure 401 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /income-from-timesheet/export/peak [post]
func (h *HttpHandler) PostExportPeak(c echo.Context) error {
	req, response, ok := h.authorizedExportPeriod(c)
	if !ok {
		return response
	}
	filename, err := h.ExportIncomeFromTimesheetUsecase.ExportPeakByStartDateAndEndDate(req.role, req.startDate, req.endDate)
	if err != nil {
		log.Println(err.Error())
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// PostExportSAP godoc
// @Summary Export individual income_from_timesheet data to a SAP file for a given period
// @Tags income-from-timesheet
// @Accept json
// @Produce json
// @Param body body models.ExportInComeSAPReq true "Export period and effective date"
// @Success 200 {array} string
// @Failure 401 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /income-from-timesheet/export/format/SAP [post]
func (h *HttpHandler) PostExportSAP(c echo.Context) error {
	if allowed, message := income.IsIncomeExportAllowed(c); !allowed {
		return c.JSON(http.StatusUnauthorized, message)
	}

	req := c.Request()
	defer req.Body.Close()

	var exReq models.ExportInComeSAPReq
	if err := json.NewDecoder(req.Body).Decode(&exReq); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	startDate, endDate, dateEff, err := exReq.ParseDates()
	if err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	filename, err := h.ExportIncomeFromTimesheetUsecase.ExportIncomeSAPByStartDateAndEndDate(exReq.Role, startDate, endDate, dateEff)
	if err != nil {
		log.Println(err.Error())
		return utils.NewError(c, http.StatusInternalServerError, errors.New("internal Server Error"))
	}

	return c.Attachment(filename, filename)
}

// authorizedMonthParam rejects callers who may not export and requests without a month.
// When ok is false the response has already been written and should be returned as-is;
// it cannot be signalled through the error alone, since echo's writers return nil on success.
func (h *HttpHandler) authorizedMonthParam(c echo.Context) (month string, response error, ok bool) {
	if allowed, message := income.IsIncomeExportAllowed(c); !allowed {
		return "", c.JSON(http.StatusUnauthorized, message), false
	}
	month = c.Param("month")
	if month == "" {
		return "", utils.NewError(c, http.StatusBadRequest, errors.New("invalid path")), false
	}
	return month, nil, true
}

type exportPeriod struct {
	role      string
	startDate time.Time
	endDate   time.Time
}

// authorizedExportPeriod rejects callers who may not export and decodes the "MM/YYYY"
// period body shared by the CSV and PEAK period exports. When ok is false the response has
// already been written and should be returned as-is.
func (h *HttpHandler) authorizedExportPeriod(c echo.Context) (period exportPeriod, response error, ok bool) {
	if allowed, message := income.IsIncomeExportAllowed(c); !allowed {
		return exportPeriod{}, c.JSON(http.StatusUnauthorized, message), false
	}

	req := c.Request()
	defer req.Body.Close()

	var t models.ExportInComeReq
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return exportPeriod{}, utils.NewError(c, http.StatusBadRequest, err), false
	}

	startDate, err := time.Parse("01/2006", t.StartDate)
	if err != nil {
		return exportPeriod{}, utils.NewError(c, http.StatusBadRequest, err), false
	}
	endDate, err := time.Parse("01/2006", t.EndDate)
	if err != nil {
		return exportPeriod{}, utils.NewError(c, http.StatusBadRequest, err), false
	}

	return exportPeriod{role: t.Role, startDate: startDate, endDate: endDate.AddDate(0, 1, 0)}, nil, true
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	incomeReader := repositories.NewIncomeFromTimesheetReader(session)
	incomeWriter := repositories.NewIncomeWriter(session)
	studentLoanRepo := repositories.NewStudentLoanRepository(session)
	sapExportFailureRepo := repositories.NewSAPExportFailureRepository(session)
	userRepo := repositories.NewUserRepository(session)
	siteRepo := site.NewRepository(session)

	ex := usecases.NewExportIncomeUsecase(
		usecases.NewIncomeFromTimesheetSource(incomeReader),
		incomeWriter,
		sapExportFailureRepo,
		file.NewCSVWriter(),
		file.NewSAPWriter(),
		file.NewPeakCSVWriter(),
		studentLoanRepo,
		userRepo,
		siteRepo,
	)
	handler := &HttpHandler{ex}

	r = r.Group("/income-from-timesheet")
	r.GET("/export/individual/:month", handler.GetExportIndividual)
	r.GET("/export/peak/individual/:month", handler.GetExportPeakIndividual)
	r.POST("/export", handler.PostExport)
	r.POST("/export/peak", handler.PostExportPeak)
	r.POST("/export/format/SAP", handler.PostExportSAP)
}
