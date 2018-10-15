package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

var (
	mockUser = models.User{
		FullName:          "นายทดสอบชอบลงทุน",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		TotalIncome:       "123123123",
		SubmitDate:        "12/12/2561",
		ThaiCitizenID:     "1234567890123",
	}

	userByte, _ = json.Marshal(mockUser)
	userJson    = string(userByte)
)

func TestCreateUser(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("CreateUser", mock.AnythingOfType("*models.User")).Return(&mockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockListUser := make([]*models.User, 0)
	mockListUser = append(mockListUser, &mockUser)

	mockUsecase.On("GetUser").Return(mockListUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.GetUser(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetUserByID", mock.AnythingOfType("string")).Return(&mockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.GetUserByID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(&mockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.UpdateUser(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("DeleteUser", mock.AnythingOfType("string")).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.DeleteUser(c)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUsecase.AssertExpectations(t)
}
