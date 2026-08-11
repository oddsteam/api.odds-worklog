package income_from_timesheet

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gitlab.odds.team/worklog/api.odds-worklog/repositories"
)

type HttpHandler struct {
	ExportIncomeFromTimesheetUsecase usecases.ForUsingExportIncomeFromTimesheet
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
	if allowed, message := income.IsIncomeExportAllowed(c); !allowed {
		return c.JSON(http.StatusUnauthorized, message)
	}
	month := c.Param("month")
	if month == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}
	filename, err := h.ExportIncomeFromTimesheetUsecase.ExportIncomeFromTimesheet("individual", month)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	incomeReader := repositories.NewIncomeFromTimesheetReader(session)
	studentLoanRepo := repositories.NewStudentLoanRepository(session)
	userRepo := user.NewRepository(session)
	siteRepo := site.NewRepository(session)
	ex := usecases.NewExportIncomeFromTimesheetUsecase(incomeReader, studentLoanRepo, userRepo, siteRepo, file.NewCSVWriter())
	handler := &HttpHandler{ex}

	r = r.Group("/income-from-timesheet")
	r.GET("/export/individual/:month", handler.GetExportIndividual)
}
