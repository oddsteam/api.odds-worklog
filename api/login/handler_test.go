package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock "gitlab.odds.team/worklog/api.odds-worklog/api/login/mock"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestLoginGoogle(t *testing.T) {
	t.Run("when request body currect then got status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().ManageLogin(mock.Login.Token).Return(&mock.MockToken, nil)

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

	})
}
