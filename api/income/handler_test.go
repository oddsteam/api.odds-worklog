package income

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
	userMocks "gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestAddIncome(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("AddIncome", &mocks.MockIncomeReq, &userMocks.MockUser).Return(&mocks.MockIncome, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.MockIncomeReqJson))
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
	handler := HttpHandler{
		Usecase: mockUsecase,
	}
	handler.AddIncome(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestUpdateIncome(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("UpdateIncome", mocks.MockIncome.ID.Hex(), &mocks.MockIncomeReq, &userMocks.MockUser).Return(&mocks.MockIncome, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.MockIncomeReqJson))
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
	c.SetParamValues(mocks.MockIncome.ID.Hex())

	handler := HttpHandler{
		Usecase: mockUsecase,
	}
	handler.UpdateIncome(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}
func TestGetIncomeStatusList(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockListUser := make([]*models.IncomeStatus, 0)
	mockListUser = append(mockListUser, &mocks.MockIncomeStatus)

	mockUsecase.On("GetIncomeStatusList").Return(mockListUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := HttpHandler{
		Usecase: mockUsecase,
	}
	handler.GetIncomeStatusList(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetIncomeByUserIdAndCurrentMonth", mocks.MockIncome.UserID).Return(&mocks.MockIncome, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(mocks.MockIncome.UserID)
	handler := HttpHandler{
		Usecase: mockUsecase,
	}
	handler.GetIncomeByUserIdAndCurrentMonth(c)

	incomeByte, _ := json.Marshal(mocks.MockIncome)
	incomeJson := string(incomeByte)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, incomeJson, rec.Body.String())
	mockUsecase.AssertExpectations(t)
}
