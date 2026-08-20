package income_from_timesheet

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

func TestGetExportIndividual(t *testing.T) {
	t.Run("when export income from timesheet succeeds it should return status OK", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("month")
		c.SetParamValues("0")

		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		startDate, endDate := models.GetStartDateAndEndDate(time.Now())
		mockRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return(mockIncomeFromTimesheetList(), nil)

		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("month")
		c.SetParamValues("0")

		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestPostExport(t *testing.T) {
	t.Run("when export income from timesheet by period succeeds it should return status OK", func(t *testing.T) {
		body := models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"}
		jsonBody, _ := json.Marshal(body)
		startDate, _ := time.Parse("01/2006", body.StartDate)
		endDate, _ := time.Parse("01/2006", body.EndDate)
		endDate = endDate.AddDate(0, 1, 0)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.EXPECT().GetAllByRoleStartDateAndEndDate(body.Role, startDate, endDate).
			Return(mockIncomeFromTimesheetList(), nil)

		handler.PostExport(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		body := models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"}
		jsonBody, _ := json.Marshal(body)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.PostExport(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when the request body is not an ExportInComeReq it should return status Bad Request", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", bytes.NewReader([]byte("not json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.PostExport(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when the export fails it should return status Internal Server Error", func(t *testing.T) {
		body := models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"}
		jsonBody, _ := json.Marshal(body)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.EXPECT().GetAllByRoleStartDateAndEndDate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		handler.PostExport(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func mockIncomeFromTimesheetList() []*models.IncomeFromTimesheet {
	record := &models.IncomeFromTimesheet{Income: models.MockIncome}
	record.UserID = bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539").Hex()
	return []*models.IncomeFromTimesheet{record}
}

func createHandlerWithMockUsecase(t *testing.T) (*HttpHandler, *gomock.Controller, *mock_usecases.MockForGettingIncomeFromTimesheetInTheMonth) {
	export, ctrl, mockRepo := usecases.CreateExportIncomeFromTimesheetUsecaseWithMock(t)
	return &HttpHandler{export}, ctrl, mockRepo
}
