package timesheet_event_log

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gitlab.odds.team/worklog/api.odds-worklog/repositories"
)

type HttpHandler struct {
	uc usecases.ForListingTimesheetEventLogs
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	lister := repositories.NewTimesheetEventLogLister(session)
	uc := usecases.NewViewTimesheetEventLogsUsecase(lister)
	handler := &HttpHandler{uc: uc}

	g := r.Group("/timesheet-event-logs")
	g.GET("", handler.List)
}

// List returns the most recently received timesheet events. Any signed-in user may read it —
// the transactional inbox is shown to the whole team.
func (h *HttpHandler) List(c echo.Context) error {
	limit := 0
	if q := c.QueryParam("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			limit = n
		}
	}

	logs, err := h.uc.List(limit)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, toResponses(logs))
}
