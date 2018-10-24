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

func TestCreateUser(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("CreateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.UserJson))
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
	mockListUser = append(mockListUser, &mocks.MockUser)

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
	mockUsecase.On("GetUserByID", mock.AnythingOfType("string")).Return(&mocks.MockUser, nil)

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
	mockUsecase.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mocks.UserJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("5bc89e26f37e2f0df54e6fef")
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

func TestLogin(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetUserByID", mock.AnythingOfType("string")).Return(&mocks.MockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.LoginJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.Login(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdatePartialUser(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetUserByID", mock.AnythingOfType("string")).Return(&mocks.MockUser, nil)
	mockUsecase.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUser, nil)
	mockIoReader := `{"fullname" : "ODDS junk","email" : "xx@c.com"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(mockIoReader))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("5bc89e26f37e2f0df54e6fef")

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.UpdatePartialUser(c)

	userByte, _ := json.Marshal(mocks.MockUser)
	UserJson := string(userByte)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, UserJson, rec.Body.String())
}

func TestUpdatePartialUserShouldReturnInternalErrorIfNoHaveRequestBody(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetUserByID", mock.AnythingOfType("string")).Return(&mocks.MockUser, nil)
	mockUsecase.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(&mocks.MockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("5bc89e26f37e2f0df54e6fef")

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.UpdatePartialUser(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
