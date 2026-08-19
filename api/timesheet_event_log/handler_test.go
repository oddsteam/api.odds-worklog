package timesheet_event_log

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestHttpHandler_List(t *testing.T) {
	newLog := func() *models.TimesheetEventLog {
		return &models.TimesheetEventLog{
			EventType: "timesheet.monthly_summary.published",
			Year:      2026,
			Month:     7,
			SummaryAt: time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
			Employee: models.TimesheetEmployee{
				Email:       "somchai@odds.team",
				EnglishName: "Somchai",
			},
			Sites: []models.TimesheetSiteSummary{
				{ClientSite: "SCB", CustomerName: "SCB Bank", WorkingDays: 18, OvertimeDays: 2},
			},
			ReceivedAt: time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC),
		}
	}

	callList := func(t *testing.T, uc *mock_usecases.MockForListingTimesheetEventLogs) *httptest.ResponseRecorder {
		t.Helper()
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &HttpHandler{uc: uc}
		assert.NoError(t, h.List(c))
		return rec
	}

	t.Run("it should expose nested site fields in camelCase so the web app can read them", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := mock_usecases.NewMockForListingTimesheetEventLogs(ctrl)
		mockUC.EXPECT().List(0).Return([]*models.TimesheetEventLog{newLog()}, nil)

		rec := callList(t, mockUC)
		assert.Equal(t, http.StatusOK, rec.Code)

		var body []map[string]interface{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		assert.Len(t, body, 1)

		sites := body[0]["sites"].([]interface{})
		site := sites[0].(map[string]interface{})
		assert.Equal(t, float64(18), site["workingDays"])
		assert.Equal(t, float64(2), site["overtimeDays"])
		assert.Equal(t, "SCB", site["clientSite"])
		assert.Equal(t, "SCB Bank", site["customerName"])

		employee := body[0]["employee"].(map[string]interface{})
		assert.Equal(t, "Somchai", employee["englishName"])
		assert.Equal(t, "somchai@odds.team", employee["email"])
	})

	t.Run("it should keep the top level fields the web app already reads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := mock_usecases.NewMockForListingTimesheetEventLogs(ctrl)
		mockUC.EXPECT().List(0).Return([]*models.TimesheetEventLog{newLog()}, nil)

		rec := callList(t, mockUC)

		var body []map[string]interface{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		assert.Equal(t, "timesheet.monthly_summary.published", body[0]["eventType"])
		assert.Equal(t, float64(2026), body[0]["year"])
		assert.Equal(t, float64(7), body[0]["month"])
		assert.Equal(t, "2026-07-31T00:00:00Z", body[0]["summaryAt"])
		assert.Equal(t, "2026-08-01T09:00:00Z", body[0]["receivedAt"])
	})

	t.Run("it should pass the limit query param to the usecase", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := mock_usecases.NewMockForListingTimesheetEventLogs(ctrl)
		mockUC.EXPECT().List(25).Return([]*models.TimesheetEventLog{}, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/?limit=25", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &HttpHandler{uc: mockUC}
		assert.NoError(t, h.List(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("it should return an empty array rather than null when there are no events", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := mock_usecases.NewMockForListingTimesheetEventLogs(ctrl)
		mockUC.EXPECT().List(0).Return(nil, nil)

		rec := callList(t, mockUC)
		assert.Equal(t, "[]", rec.Body.String())
	})
}
