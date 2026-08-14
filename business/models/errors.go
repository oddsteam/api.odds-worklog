package models

import (
	"errors"
	"log"
)

var (
	ErrInvalidUserRole    = errors.New("Invalid user role")
	ErrInvalidUserVat     = errors.New("Invalid user vat.")
	ErrConflict           = errors.New("Item already exist")
	ErrPermissionDenied   = errors.New("Permission denied.")
	ErrInvalidFormat      = errors.New("Invalid format")
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
