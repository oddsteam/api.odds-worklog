package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestConfigError(t *testing.T) {
	assert.Equal(t, "Item not found", ErrNotFound.Error())
	assert.Equal(t, "Cannot be Deleted", ErrCannotBeDeleted.Error())
	assert.Equal(t, "Item already exist", ErrConflict.Error())
	assert.Equal(t, "Invalid path", ErrInvalidPath.Error())
	assert.Equal(t, "Invalid format", ErrInvalidFormat.Error())
	assert.Equal(t, "Invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "Bad request", ErrBadRequest.Error())
	assert.Equal(t, "Invalid flag", ErrInvalidFlag.Error())
	assert.Equal(t, "Email is not account @odds.team", ErrEmailIsNotOddsTeam.Error())
	assert.Equal(t, "Token is not account @odds.team", ErrTokenIsNotOddsTeam.Error())
	assert.Equal(t, "Invalid user role", ErrInvalidUserRole.Error())
	assert.Equal(t, "Save transcript failed", ErrSaveTranscript.Error())
	assert.Equal(t, "Not PDF file", ErrNotPDFFile.Error())
	assert.Equal(t, "No transcript file.", ErrNoTranscriptFile.Error())
	assert.Equal(t, "No degree certificate file.", ErrNoDegreeCertificateFile.Error())
	assert.Equal(t, "Permission denied.", ErrPermissionDenied.Error())
	assert.Equal(t, "Invalid user vat.", ErrInvalidUserVat.Error())
	assert.Equal(t, "Not image file", ErrNotImageFile.Error())
	assert.Equal(t, "Invalid amount.", ErrInvalidAmount.Error())
}

func TestNewError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	NewError(c, http.StatusOK, errors.New("ok"))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `{"code":200,"message":"ok"}`, rec.Body.String())
}
