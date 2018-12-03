package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	t.Run("when email is valid, then return nil", func(t *testing.T) {
		assert.Nil(t, ValidateEmail("abc@mail.com"))
	})

	t.Run("when email is invalid, then return ErrInvalidFormat", func(t *testing.T) {
		assert.EqualError(t, ValidateEmail(""), ErrInvalidFormat.Error())
		assert.EqualError(t, ValidateEmail("a"), ErrInvalidFormat.Error())
		assert.EqualError(t, ValidateEmail("@mail.com"), ErrInvalidFormat.Error())
	})
}
