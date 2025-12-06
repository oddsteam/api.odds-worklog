package income

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/entity"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/entity/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/api/usecases"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestAddIncome(t *testing.T) {
	t.Run("when request body isRequestValid then got status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().AddIncome(&entity.MockIncomeReq, userMock.User.ID.Hex()).Return(&entity.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(entity.MockIncomeReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.AddIncome(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when add income but request body is not IncomeReq it should be return status Unprocessable Entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(entity.MockIncomeResJson))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues(entity.MockIncome.ID.Hex())

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.AddIncome(c)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

}

func TestUpdateIncome(t *testing.T) {
	t.Run("when update income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().UpdateIncome(entity.MockIncome.ID.Hex(), &entity.MockIncomeReq, userMock.User.ID.Hex()).Return(&entity.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(entity.MockIncomeReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues(entity.MockIncome.ID.Hex())

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
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

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.UpdateIncome(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update income but request body is not IncomeReq it should be return status Unprocessable Entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(entity.MockIncomeResJson))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues(entity.MockIncome.ID.Hex())

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
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
		mockListUser = append(mockListUser, &entity.MockCorporateIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("corporate", true).Return(mockListUser, nil)

		e := echo.New()

		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetCorporateIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get corporate income status list it should not return individual list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &entity.MockIndividualIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("corporate", true).Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetCorporateIncomeStatus(c)
		incomeByte, _ := json.Marshal(entity.MockIndividualIncomeStatus)
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
		mockListUser = append(mockListUser, &entity.MockIndividualIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("individual", true).Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetIndividualIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get individual income status list it should not return corporate list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &entity.MockCorporateIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("individual", true).Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetIndividualIncomeStatus(c)
		incomeByte, _ := json.Marshal(entity.MockCorporateIncomeStatus)
		incomeJson := string(incomeByte)
		assert.NotEqual(t, incomeJson, rec.Body)

	})
}

func TestGetIncomeGetIncomeCurrentMonthByUserId(t *testing.T) {
	t.Run("when get income by user id in current month success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetIncomeByUserIdAndCurrentMonth(entity.MockIncome.UserID).Return(&entity.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetIncomeCurrentMonthByUserId(c)

		incomeByte, _ := json.Marshal(entity.MockIncome)
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

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetIncomeCurrentMonthByUserId(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
func TestGetIncomeGetIncomeAllMonthByUserId(t *testing.T) {
	t.Run("when get income by user id in all month success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetIncomeByUserIdAllMonth(entity.MockIncome.UserID).Return(entity.MockIncomeList, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler, ctrl := createHandlerWithMockUsecases(t, mockUsecase)
		defer ctrl.Finish()
		handler.GetIncomeAllMonthByUserId(c)

		incomeByte, _ := json.Marshal(entity.MockIncomeList)
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

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("month")
		c.SetParamValues("1")
		handler, ctrl, mockRepo := createHandlerWithMockUsecasesAndRepo(t, mockUsecase)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllIncomeOfPreviousMonthByRole(entity.MockIncomeList)
		handler.GetExportCorporate(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetExportIndividualIncomeStatus(t *testing.T) {
	t.Run("when export individual income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("month")
		c.SetParamValues("1")
		handler, ctrl, mockRepo := createHandlerWithMockUsecasesAndRepo(t, mockUsecase)
		defer ctrl.Finish()
		mockRepo.ExpectGetAllIncomeOfPreviousMonthByRole(entity.MockIncomeList)
		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

type ExportInComeSAPReq struct {
	Role          string
	DateEffective string
	StartDate     string
	EndDate       string
}

func TestPostExportSAPIncome(t *testing.T) {
	t.Run("when export corporate income as SAP format by period time success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		body := ExportInComeSAPReq{
			Role:          "corporate",
			DateEffective: "30/09/2025",
			StartDate:     "09/2025",
			EndDate:       "10/2025",
		}
		jsonBody, _ := json.Marshal(body)
		startDate, _ := time.Parse("01/2006", body.StartDate)
		endDate, _ := time.Parse("01/2006", body.EndDate)
		endDate = endDate.AddDate(0, 1, 0)

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler, ctrl, mockRepo := createHandlerWithMockUsecasesAndRepo(t, mockUsecase)
		defer ctrl.Finish()
		mockRepo.GetAllIncomeByRoleStartDateAndEndDate(entity.MockIncomeList, body.Role, startDate, endDate)
		handler.PostExportSAP(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when export individual income as SAP format by period time success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		body := ExportInComeSAPReq{
			Role:          "individual",
			DateEffective: "30/09/2025",
			StartDate:     "09/2025",
			EndDate:       "10/2025",
		}
		jsonBody, _ := json.Marshal(body)
		startDate, _ := time.Parse("01/2006", body.StartDate)
		endDate, _ := time.Parse("01/2006", body.EndDate)
		endDate = endDate.AddDate(0, 1, 0)

		mockUsecase := incomeMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler, ctrl, mockRepo := createHandlerWithMockUsecasesAndRepo(t, mockUsecase)
		defer ctrl.Finish()
		mockRepo.GetAllIncomeByRoleStartDateAndEndDate(entity.MockIncomeList, body.Role, startDate, endDate)
		handler.PostExportSAP(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func createHandlerWithMockUsecases(t *testing.T, mockUsecase *incomeMock.MockUsecase) (*HttpHandler, *gomock.Controller) {
	export, ctrl, _ := usecases.CreateExportIncomeUsecaseWithMock(t)
	return &HttpHandler{mockUsecase, export}, ctrl
}

func createHandlerWithMockUsecasesAndRepo(t *testing.T, mockUsecase *incomeMock.MockUsecase) (*HttpHandler, *gomock.Controller, *usecases.MockIncomeRepository) {
	export, ctrl, mockRepo := usecases.CreateExportIncomeUsecaseWithMock(t)
	return &HttpHandler{mockUsecase, export}, ctrl, mockRepo
}
