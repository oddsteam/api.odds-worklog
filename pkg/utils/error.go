package utils

import (
	"errors"

	"github.com/labstack/echo"
)

var (
	ErrNotFound                = errors.New("Item not found")
	ErrCannotBeDeleted         = errors.New("Cannot be Deleted")
	ErrConflict                = errors.New("Item already exist")
	ErrInvalidPath             = errors.New("Invalid path")
	ErrInvalidFormat           = errors.New("Invalid format")
	ErrInvalidToken            = errors.New("Invalid token")
	ErrBadRequest              = errors.New("Bad request")
	ErrInvalidFlag             = errors.New("Invalid flag")
	ErrEmailIsNotOddsTeam      = errors.New("Email is not account @odds.team")
	ErrTokenIsNotOddsTeam      = errors.New("Token is not account @odds.team")
	ErrInvalidUserRole         = errors.New("Invalid user role")
	ErrSaveTranscript          = errors.New("Save transcript failed")
	ErrNotPDFFile              = errors.New("Not PDF file")
	ErrNoTranscriptFile        = errors.New("No transcript file.")
	ErrNoImageProfileFile      = errors.New("No image profile.")
	ErrNoDegreeCertificateFile = errors.New("No degree certificate file.")
	ErrNoIDCardFile            = errors.New("No idcard file.")
	ErrPermissionDenied        = errors.New("Permission denied.")
	ErrInvalidUserVat          = errors.New("Invalid user vat.")
	ErrCustomerNotFound        = errors.New("Customer is not found.")
	ErrNotImageFile            = errors.New("Not image file")
	ErrInvalidAmount           = errors.New("Invalid amount.")
)

func NewError(ctx echo.Context, status int, err error) error {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	return ctx.JSON(status, er)
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
