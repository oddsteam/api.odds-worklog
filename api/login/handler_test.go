package login

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	mock "gitlab.odds.team/worklog/api.odds-worklog/api/login/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestLoginGoogle(t *testing.T) {
	t.Run("when request body currect then got status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "abc@mail.com"
		user := new(models.User)
		user.Email = email

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetTokenInfo(mock.Login.Token).Return(&mock.MockTokenInfo, nil)
		mockUsecase.EXPECT().CreateUserAndValidateEmail(mock.MockTokenInfo.Email).Return(user, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.LoginJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.loginGoogle(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when request body wrong then got status unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.loginGoogle(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, `{"code":401,"message":"Bad request"}`, rec.Body.String())
	})

	t.Run("when request token is empty then got status unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loginJSON := `{"token": ""}`
		mockUsecase := mock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.loginGoogle(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, `{"code":401,"message":"Bad request"}`, rec.Body.String())
	})

	t.Run("when request token is invalid then got status unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loginJSON := `{"token": "1234"}`

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetTokenInfo("1234").Return(nil, utils.ErrInvalidToken)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.loginGoogle(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when create user failed then got status unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "abc@mail.com"
		user := new(models.User)
		user.Email = email

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetTokenInfo(mock.Login.Token).Return(&mock.MockTokenInfo, nil)
		mockUsecase.EXPECT().CreateUserAndValidateEmail(mock.MockTokenInfo.Email).Return(user, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.LoginJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.loginGoogle(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
