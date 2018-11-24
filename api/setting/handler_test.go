package setting_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/api/setting"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	"github.com/labstack/echo"
)

type MockRepositorySuccess struct{}

func NewMockRepositorySuccess() setting.Repository {
	return MockRepositorySuccess{}
}

func (fs MockRepositorySuccess) GetReminder() (*models.Reminder, error) {
	reminder := new(models.Reminder)
	reminder.Name = "reminder"
	return reminder, nil
}

func (fs MockRepositorySuccess) SaveReminder(reminder *models.Reminder) (*models.Reminder, error) {
	return reminder, nil
}

type MockRepositoryFail struct{}

func NewMockRepositoryFail() setting.Repository {
	return MockRepositoryFail{}
}

func (fs MockRepositoryFail) GetReminder() (*models.Reminder, error) {
	return nil, errors.New("Data not found in DB")
}

func (fs MockRepositoryFail) SaveReminder(reminder *models.Reminder) (*models.Reminder, error) {
	return nil, errors.New("Can not save data in DB")
}
func TestGetReminderSettingSuccess(t *testing.T) {
	// use the mocked object
	mockRepository := NewMockRepositorySuccess()
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setting.Get(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetReminderSettingFail(t *testing.T) {
	// use the mocked object
	mockRepository := NewMockRepositoryFail()
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setting.Get(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSaveReminderSuccess(t *testing.T) {
	reader := strings.NewReader("{}")

	mockRepository := NewMockRepositorySuccess()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setting.Save(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSaveReminderFail(t *testing.T) {
	reader := strings.NewReader("{}")

	mockRepository := NewMockRepositoryFail()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setting.Save(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
