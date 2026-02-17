package models

import (
	"errors"
	"log"
)

var (
	ErrInvalidUserRole = errors.New("Invalid user role")
	ErrInvalidUserVat  = errors.New("Invalid user vat.")
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
