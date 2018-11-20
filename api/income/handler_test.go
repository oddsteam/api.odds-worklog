package income

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	jwt "github.com/dgrijalva/jwt-go"
// 	"github.com/labstack/echo"
// 	"github.com/stretchr/testify/assert"
// 	"gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
// 	userMocks "gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
// 	"gitlab.odds.team/worklog/api.odds-worklog/models"
// )

// func TestAddIncome(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockUsecase.On("AddIncome", &mocks.MockIncomeReq, &userMocks.MockUser).Return(&mocks.MockIncome, nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.MockIncomeReqJson))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 	claims := &models.JwtCustomClaims{
// 		&userMocks.MockUser,
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.Set("user", token)
// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.AddIncome(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	mockUsecase.AssertExpectations(t)
// }

// func TestUpdateIncome(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockUsecase.On("UpdateIncome", mocks.MockIncome.ID.Hex(), &mocks.MockIncomeReq, &userMocks.MockUser).Return(&mocks.MockIncome, nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.MockIncomeReqJson))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 	claims := &models.JwtCustomClaims{
// 		&userMocks.MockUser,
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.Set("user", token)
// 	c.SetParamNames("id")
// 	c.SetParamValues(mocks.MockIncome.ID.Hex())

// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.UpdateIncome(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	mockUsecase.AssertExpectations(t)
// }

// func TestGetCorporateIncomeStatus(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockListUser := make([]*models.IncomeStatus, 0)
// 	mockListUser = append(mockListUser, &mocks.MockIncomeStatus)

// 	mockUsecase.On("GetIncomeStatusList", "Y").Return(mockListUser, nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.GetCorporateIncomeStatus(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	mockUsecase.AssertExpectations(t)
// }

// func TestGetIndividualIncomeStatus(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockListUser := make([]*models.IncomeStatus, 0)
// 	mockListUser = append(mockListUser, &mocks.MockIncomeStatus)

// 	mockUsecase.On("GetIncomeStatusList", "N").Return(mockListUser, nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.GetIndividualIncomeStatus(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	b, _ := json.Marshal(mockListUser)
// 	assert.Equal(t, string(b), rec.Body.String())
// 	mockUsecase.AssertExpectations(t)
// }

// func TestGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockUsecase.On("GetIncomeByUserIdAndCurrentMonth", mocks.MockIncome.UserID).Return(&mocks.MockIncome, nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues(mocks.MockIncome.UserID)
// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.GetIncomeByUserIdAndCurrentMonth(c)

// 	incomeByte, _ := json.Marshal(mocks.MockIncome)
// 	incomeJson := string(incomeByte)
// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	assert.Equal(t, incomeJson, rec.Body.String())
// 	mockUsecase.AssertExpectations(t)
// }

// func TestGetExportCorporateIncomeStatus(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockUsecase.On("ExportIncome", "Y").Return("test.csv", nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.GetExportCorporate(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	mockUsecase.AssertExpectations(t)
// }

// func TestGetExportIndividualIncomeStatus(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockUsecase.On("ExportIncome", "N").Return("test.csv", nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.GetExportIndividual(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	mockUsecase.AssertExpectations(t)
// }

// func TestDropIncome(t *testing.T) {
// 	mockUsecase := new(mocks.Usecase)
// 	mockUsecase.On("DropIncome").Return(nil)

// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	handler := HttpHandler{
// 		Usecase: mockUsecase,
// 	}
// 	handler.DropIncome(c)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	assert.Equal(t, "{\"message\":\"DropIncome Success!\"}", rec.Body.String())
// 	mockUsecase.AssertExpectations(t)
// }

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockIncome "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMocks "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestAddIncome(t *testing.T) {
	t.Run("when request body isRequestValid then got status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().AddIncome(&mockIncome.MockIncomeReq, &userMocks.MockUser).Return(&mockIncome.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mockIncome.MockIncomeReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		claims := &models.JwtCustomClaims{
			&userMocks.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)

		handler := &HttpHandler{mockUsecase}
		handler.AddIncome(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when add income but request body is not IncomeReq it should be return status Unprocessable Entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mockIncome.MockIncomeResJson))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &models.JwtCustomClaims{
			&userMocks.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(mockIncome.MockIncome.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.AddIncome(c)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

}

func TestUpdateIncome(t *testing.T) {
	t.Run("when update income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().UpdateIncome(mockIncome.MockIncome.ID.Hex(), &mockIncome.MockIncomeReq, &userMocks.MockUser).Return(&mockIncome.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mockIncome.MockIncomeReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		claims := &models.JwtCustomClaims{
			&userMocks.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(mockIncome.MockIncome.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.UpdateIncome(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when update income but no have id it should be return status Bad Request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)

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

		mockUsecase := mockIncome.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mockIncome.MockIncomeResJson))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &models.JwtCustomClaims{
			&userMocks.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(mockIncome.MockIncome.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.UpdateIncome(c)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

}

func TestGetCorporateIncomeStatus(t *testing.T) {
	t.Run("when get corporate income status list is success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &mockIncome.MockCorporateIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("Y").Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetCorporateIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get corporate income status list it should not return individual list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &mockIncome.MockIndividualIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("Y").Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetCorporateIncomeStatus(c)
		incomeByte, _ := json.Marshal(mockIncome.MockIndividualIncomeStatus)
		incomeJson := string(incomeByte)
		assert.NotEqual(t, incomeJson, rec.Body)

	})
}

func TestGetIndividualIncomeStatus(t *testing.T) {
	t.Run("when get individual income status list is success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &mockIncome.MockIndividualIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("N").Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetIndividualIncomeStatus(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get individual income status list it should not return corporate list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockListUser := make([]*models.IncomeStatus, 0)
		mockListUser = append(mockListUser, &mockIncome.MockCorporateIncomeStatus)
		mockUsecase.EXPECT().GetIncomeStatusList("N").Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetIndividualIncomeStatus(c)
		incomeByte, _ := json.Marshal(mockIncome.MockCorporateIncomeStatus)
		incomeJson := string(incomeByte)
		assert.NotEqual(t, incomeJson, rec.Body)

	})
}

func TestGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
	t.Run("when get income by user id in current month success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetIncomeByUserIdAndCurrentMonth(mockIncome.MockIncome.UserID).Return(&mockIncome.MockIncome, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(mockIncome.MockIncome.UserID)

		handler := &HttpHandler{mockUsecase}
		handler.GetIncomeByUserIdAndCurrentMonth(c)

		incomeByte, _ := json.Marshal(mockIncome.MockIncome)
		incomeJson := string(incomeByte)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, incomeJson, rec.Body.String())
	})

	t.Run("when get income by user id in current month is no have id it should be return status Bad Request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetIncomeByUserIdAndCurrentMonth(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetExportCorporateIncomeStatus(t *testing.T) {
	t.Run("when export corporate income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ExportIncome("Y").Return("test.csv", nil)
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetExportCorporate(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetExportIndividualIncomeStatus(t *testing.T) {
	t.Run("when export individual income success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ExportIncome("N").Return("test.csv", nil)
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetExportIndividual(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestDropIncome(t *testing.T) {
	t.Run("when drop table success it should be return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mockIncome.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().DropIncome().Return(nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.DropIncome(c)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"message\":\"DropIncome Success!\"}", rec.Body.String())

	})
}
