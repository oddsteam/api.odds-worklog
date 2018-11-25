package reminder_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/api/reminder"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	"github.com/labstack/echo"
)

type MockRepositorySuccess struct{}

func NewMockRepositorySuccess() reminder.Repository {
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

func NewMockRepositoryFail() reminder.Repository {
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

	reminder.GetReminder(c, mockRepository)
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

	reminder.GetReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSaveReminderSuccess(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	request.Setting.Date = "25"
	request.Setting.Message = "TEST"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)

	mockRepository := NewMockRepositorySuccess()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSaveReminderShouldInternalServerErr_WhenCanNotSaveIntoDB(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	request.Setting.Date = "25"
	request.Setting.Message = "TEST"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositoryFail()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSaveReminderShouldBadRequest_WhenRequestIsEmpty(t *testing.T) {
	mockRepository := NewMockRepositoryFail()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSaveReminderShouldBadRequest_WhenRequestNameIsEmpty(t *testing.T) {
	request := new(models.Reminder)
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositoryFail()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSaveReminderShouldBadRequest_WhenRequestSettingDateIsEmpty(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositoryFail()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSaveReminderShouldBadRequest_WhenRequestSettingMessageIsEmpty(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	request.Setting.Date = "25"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositoryFail()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSaveReminderShouldSuccess_WhenRequestSettingDateIs26(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	request.Setting.Date = "26"
	request.Setting.Message = "TEST"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositorySuccess()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSaveReminderShouldSuccess_WhenRequestSettingDateIs27(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	request.Setting.Date = "27"
	request.Setting.Message = "TEST"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositorySuccess()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSaveReminderShouldBadRequest_WhenRequestSettingDateIs1(t *testing.T) {
	request := new(models.Reminder)
	request.Name = "reminder"
	request.Setting.Date = "1"
	request.Setting.Message = "TEST"
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	mockRepository := NewMockRepositorySuccess()

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	reminder.SaveReminder(c, mockRepository)
	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
