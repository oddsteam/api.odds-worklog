package income_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
)

func TestAddIncome(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("AddIncome", mock.AnythingOfType("*models.Income")).Return(&mocks.MockIncome, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mocks.AddIncomeJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidG9tQG9kZHMudGVhbSIsImFkbWluIjp0cnVlLCJleHAiOjE1NDAyODIxMDN9.a1B6D2RDjeFmBz8RHVgaDGHLMifb5Ml9Dzz1CGvsOKo")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := income.HttpHandler{
		Usecase: mockUsecase,
	}
	handler.AddIncome(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUsecase.AssertExpectations(t)
}
