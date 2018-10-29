package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
)

func TestLogin(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetUserByID", "5bbcf2f90fd2df527bc39539").Return(&mocks.MockUserById, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.LoginJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	login(c, mockRepo)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}
