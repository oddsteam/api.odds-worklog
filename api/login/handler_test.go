package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocker "github.com/stretchr/testify/mock"
	mock "gitlab.odds.team/worklog/api.odds-worklog/api/login/mock"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
)

func TestLogin(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetUserByID", "5bbcf2f90fd2df527bc39539").Return(&mocks.MockUserById, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.LoginJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	login(c, mockRepo)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestLoginGoogle(t *testing.T) {
	t.Run("when request body currect then got status OK", func(t *testing.T) {
		mockRepo := new(mocks.Repository)
		mockRepo.On("CreateUser", mocker.Anything).Return(&mocks.MockUser, nil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogin := mock.NewMockUsecase(ctrl)
		mockLogin.EXPECT().GetTokenInfo(mock.Login.Token).Return(&mock.MockTokenInfo, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.LoginJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		loginGoogle(mockLogin, c, mockRepo)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("when request body wrong then got status unauthorized", func(t *testing.T) {

	})
}
