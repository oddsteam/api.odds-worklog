package utils

import "github.com/labstack/echo"

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
