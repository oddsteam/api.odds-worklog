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
	ucmock "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

func TestGetExportIndividual(t *testing.T) {
	t.Run("when export income from timesheet succeeds it should return status OK", func(t *testing.T) {
		c, rec := getContext(userMock.TokenAdmin, "0")
		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllInTheMonth("individual", "0", mockIncomeFromTimesheetList())

		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		c, rec := getContext(userMock.TokenUser, "0")
		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when the month param is missing it should return status Bad Request", func(t *testing.T) {
		e := echo.New()
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(echo.GET, "/", nil), rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when the export fails it should return status Internal Server Error", func(t *testing.T) {
		c, rec := getContext(userMock.TokenAdmin, "1")
		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllFails()

		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetExportPeakIndividual(t *testing.T) {
	t.Run("when export peak from timesheet succeeds it should return status OK", func(t *testing.T) {
		c, rec := getContext(userMock.TokenAdmin, "1")
		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllInTheMonth("individual", "1", mockIncomeFromTimesheetList())

		handler.GetExportPeakIndividual(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		c, rec := getContext(userMock.TokenUser, "0")
		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.GetExportPeakIndividual(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestPostExport(t *testing.T) {
	t.Run("when export income from timesheet by period succeeds it should return status OK", func(t *testing.T) {
		body := models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"}
		startDate, endDate := periodOf(body)
		c, rec := postContext(userMock.TokenAdmin, body)

		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllByPeriod(body.Role, startDate, endDate, mockIncomeFromTimesheetList())

		handler.PostExport(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		c, rec := postContext(userMock.TokenUser, models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"})
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

	t.Run("when the period is not MM/YYYY it should return status Bad Request", func(t *testing.T) {
		c, rec := postContext(userMock.TokenAdmin, models.ExportInComeReq{Role: "individual", StartDate: "2026-06", EndDate: "06/2026"})
		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.PostExport(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when the export fails it should return status Internal Server Error", func(t *testing.T) {
		c, rec := postContext(userMock.TokenAdmin, models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"})
		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllFails()

		handler.PostExport(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestPostExportPeak(t *testing.T) {
	t.Run("when export peak from timesheet by period succeeds it should return status OK", func(t *testing.T) {
		body := models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"}
		startDate, endDate := periodOf(body)
		c, rec := postContext(userMock.TokenAdmin, body)

		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllByPeriod(body.Role, startDate, endDate, mockIncomeFromTimesheetList())

		handler.PostExportPeak(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		c, rec := postContext(userMock.TokenUser, models.ExportInComeReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026"})
		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.PostExportPeak(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestPostExportSAP(t *testing.T) {
	t.Run("when export SAP from timesheet succeeds it should return status OK", func(t *testing.T) {
		body := models.ExportInComeSAPReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026", DateEffective: "01/06/2026"}
		startDate, endDate, _, err := body.ParseDates()
		assert.NoError(t, err)
		c, rec := postContext(userMock.TokenAdmin, body)

		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllByPeriod(body.Role, startDate, endDate, mockIncomeFromTimesheetList())

		handler.PostExportSAP(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when the caller is not allowed to export it should return status Unauthorized", func(t *testing.T) {
		c, rec := postContext(userMock.TokenUser, models.ExportInComeSAPReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026", DateEffective: "01/06/2026"})
		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.PostExportSAP(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when the dates cannot be parsed it should return status Bad Request", func(t *testing.T) {
		c, rec := postContext(userMock.TokenAdmin, models.ExportInComeSAPReq{Role: "individual", StartDate: "nope", EndDate: "06/2026", DateEffective: "01/06/2026"})
		handler, ctrl, _ := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()

		handler.PostExportSAP(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when the export fails it should return status Internal Server Error", func(t *testing.T) {
		c, rec := postContext(userMock.TokenAdmin, models.ExportInComeSAPReq{Role: "individual", StartDate: "06/2026", EndDate: "06/2026", DateEffective: "01/06/2026"})
		handler, ctrl, mockRepo := createHandlerWithMockUsecase(t)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllFails()

		handler.PostExportSAP(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetIndividualIncomeStatus(t *testing.T) {
	t.Run("when get individual income status list from timesheet succeeds it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockListIncomeStatus := ucmock.NewMockForUsingListIncomeStatus(ctrl)
		mockListIncomeStatus.EXPECT().GetIncomeStatusList("individual", true).
			Return([]*models.IncomeStatus{&models.MockIndividualIncomeStatus}, nil)

		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl := createHandlerWithMockListAndGetUsecases(t, mockListIncomeStatus, ucmock.NewMockForUsingGetIncome(ctrl))
		defer ctrl.Finish()
		handler.GetIndividualIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetIncomeCurrentMonthByUserId(t *testing.T) {
	t.Run("when get income from timesheet by user id in current month succeeds it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGetIncomeUsecase := ucmock.NewMockForUsingGetIncome(ctrl)
		mockGetIncomeUsecase.EXPECT().GetIncomeByCurrentMonth(models.MockIncome.UserID).Return(&models.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler, ctrl := createHandlerWithMockListAndGetUsecases(t, ucmock.NewMockForUsingListIncomeStatus(ctrl), mockGetIncomeUsecase)
		defer ctrl.Finish()
		handler.GetIncomeCurrentMonthByUserId(c)

		incomeByte, _ := json.Marshal(models.MockIncome)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(incomeByte), rec.Body.String())
	})

	t.Run("when the requested id does not match the caller it should return status Bad Request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("not-the-caller")

		handler, ctrl := createHandlerWithMockListAndGetUsecases(t, ucmock.NewMockForUsingListIncomeStatus(ctrl), ucmock.NewMockForUsingGetIncome(ctrl))
		defer ctrl.Finish()
		handler.GetIncomeCurrentMonthByUserId(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetIncomeAllMonthByUserId(t *testing.T) {
	t.Run("when get income from timesheet by user id in all month succeeds it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGetIncomeUsecase := ucmock.NewMockForUsingGetIncome(ctrl)
		mockGetIncomeUsecase.EXPECT().GetIncomeByAllMonth(models.MockIncome.UserID).Return(models.MockIncomeList, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler, ctrl := createHandlerWithMockListAndGetUsecases(t, ucmock.NewMockForUsingListIncomeStatus(ctrl), mockGetIncomeUsecase)
		defer ctrl.Finish()
		handler.GetIncomeAllMonthByUserId(c)

		incomeByte, _ := json.Marshal(models.MockIncomeList)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(incomeByte), rec.Body.String())
	})
}

func getContext(user interface{}, month string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(echo.GET, "/", nil), rec)
	c.Set("user", user)
	c.SetParamNames("month")
	c.SetParamValues(month)
	return c, rec
}

func postContext(user interface{}, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	jsonBody, _ := json.Marshal(body)
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)
	return c, rec
}

func periodOf(body models.ExportInComeReq) (time.Time, time.Time) {
	startDate, _ := time.Parse("01/2006", body.StartDate)
	endDate, _ := time.Parse("01/2006", body.EndDate)
	return startDate, endDate.AddDate(0, 1, 0)
}

func mockIncomeFromTimesheetList() []*models.IncomeFromTimesheet {
	record := &models.IncomeFromTimesheet{Income: models.MockIncome}
	record.UserID = bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539").Hex()
	return []*models.IncomeFromTimesheet{record}
}

func createHandlerWithMockUsecase(t *testing.T) (*HttpHandler, *gomock.Controller, *usecases.MockIncomeFromTimesheetRepository) {
	export, ctrl, mockRepo := usecases.CreateExportIncomeFromTimesheetUsecaseWithMock(t)
	return &HttpHandler{ExportIncomeFromTimesheetUsecase: export}, ctrl, mockRepo
}

func createHandlerWithMockListAndGetUsecases(t *testing.T, mockListIncomeStatus usecases.ForUsingListIncomeStatus, mockGetIncome usecases.ForUsingGetIncome) (*HttpHandler, *gomock.Controller) {
	export, ctrl, _ := usecases.CreateExportIncomeFromTimesheetUsecaseWithMock(t)
	return &HttpHandler{
		ListIncomeStatusUsecase:          mockListIncomeStatus,
		GetIncomeUsecase:                 mockGetIncome,
		ExportIncomeFromTimesheetUsecase: export,
	}, ctrl
}
