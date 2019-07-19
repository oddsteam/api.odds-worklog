package income

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestAddIncome(t *testing.T) {
	t.Run("when request body isRequestValid then got status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().AddIncome(&incomeMock.MockIncomeReq, &userMock.User).Return(&incomeMock.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(incomeMock.MockIncomeReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		handler := &HttpHandler{mockUsecase}
		handler.AddIncome(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when add income but request body is not IncomeReq it should be return status Unprocessable Entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(incomeMock.MockIncomeResJson))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues(incomeMock.MockIncome.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.AddIncome(c)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

}

func TestUpdateIncome(t *testing.T) {
	t.Run("when update income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().UpdateIncome(incomeMock.MockIncome.ID.Hex(), &incomeMock.MockIncomeReq, &userMock.User).Return(&incomeMock.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(incomeMock.MockIncomeReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues(incomeMock.MockIncome.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.UpdateIncome(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when update income but no have id it should be return status Bad Request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.UpdateIncome(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update income but request body is not IncomeReq it should be return status Unprocessable Entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(incomeMock.MockIncomeResJson))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues(incomeMock.MockIncome.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.UpdateIncome(c)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

}

func TestGetCorporateIncomeStatus(t *testing.T) {
	t.Run("when get corporate income status list is success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &incomeMock.MockCorporateIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("corporate", true).Return(mockListUser, nil)

		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.GetCorporateIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get corporate income status list it should not return individual list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &incomeMock.MockIndividualIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("corporate", true).Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.GetCorporateIncomeStatus(c)
		incomeByte, _ := json.Marshal(incomeMock.MockIndividualIncomeStatus)
		incomeJson := string(incomeByte)
		assert.NotEqual(t, incomeJson, rec.Body)

	})
}

func TestGetIndividualIncomeStatus(t *testing.T) {
	t.Run("when get individual income status list is success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &incomeMock.MockIndividualIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("individual", true).Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.GetIndividualIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get individual income status list it should not return corporate list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &incomeMock.MockCorporateIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("individual", true).Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.GetIndividualIncomeStatus(c)
		incomeByte, _ := json.Marshal(incomeMock.MockCorporateIncomeStatus)
		incomeJson := string(incomeByte)
		assert.NotEqual(t, incomeJson, rec.Body)

	})
}

func TestGetIncomeGetIncomeCurrentMonthByUserId(t *testing.T) {
	t.Run("when get income by user id in current month success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetIncomeByUserIdAndCurrentMonth(incomeMock.MockIncome.UserID).Return(&incomeMock.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.GetIncomeCurrentMonthByUserId(c)

		incomeByte, _ := json.Marshal(incomeMock.MockIncome)
		incomeJson := string(incomeByte)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, incomeJson, rec.Body.String())
	})

	t.Run("when get income by user id in current month is no have id it should be return status Bad Request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mockUser := userMock.User
		mockUser.ID = ""
		c.Set("user", userMock.TokenUser)

		handler := &HttpHandler{mockUsecase}
		handler.GetIncomeCurrentMonthByUserId(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
func TestGetIncomeGetIncomeAllMonthByUserId(t *testing.T) {
	t.Run("when get income by user id in all month success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetIncomeByUserIdAllMonth(incomeMock.MockIncome.UserID).Return(incomeMock.MockIncomeList, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.GetIncomeAllMonthByUserId(c)

		incomeByte, _ := json.Marshal(incomeMock.MockIncomeList)
		incomeJson := string(incomeByte)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, incomeJson, rec.Body.String())
	})
}

// func TestGetExportPdf(t *testing.T) {
// 	t.Run("when export pdf success it should be return status OK", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		mockUsecase := incomeMock.NewMockUsecase(ctrl)
// 		mockUsecase.EXPECT().ExportPdf().Return("test.pdf", nil)
// 		e := echo.New()
// 		req := httptest.NewRequest(echo.GET, "/", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		handler := &HttpHandler{mockUsecase}
// 		handler.GetExportCorporate(c)

// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	})
// }

func TestGetExportCorporateIncomeStatus(t *testing.T) {
	t.Run("when export corporate income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ExportIncome("corporate", "1").Return("test.csv", nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("month")
		c.SetParamValues("1")
		handler := &HttpHandler{mockUsecase}
		handler.GetExportCorporate(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetExportDifferentCorporate(t *testing.T) {
	t.Run("when export different corporate income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ExportIncomeNotExport("corporate").Return("test.csv", nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		handler := &HttpHandler{mockUsecase}
		handler.GetExportDifferentCorporate(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetExportIndividualIncomeStatus(t *testing.T) {
	t.Run("when export individual income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ExportIncome("individual", "1").Return("test.csv", nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("month")
		c.SetParamValues("1")
		handler := &HttpHandler{mockUsecase}
		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetExportDifferentIndividuals(t *testing.T) {
	t.Run("when export different individuals income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ExportIncomeNotExport("individual").Return("test.csv", nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		handler := &HttpHandler{mockUsecase}
		handler.GetExportDifferentIndividuals(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
