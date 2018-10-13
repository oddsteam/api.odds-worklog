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
)

func TestCreateUser(t *testing.T) {
	tempMockUser := mockUser
	mockUsecase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockUser)
	assert.NoError(t, err)

	mockUsecase.On("CreateUser", mock.AnythingOfType("*models.User")).Return(&mockUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user")

	handler := user.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUsecase.AssertExpectations(t)
}
