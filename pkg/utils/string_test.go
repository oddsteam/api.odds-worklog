package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFirstUpper(t *testing.T) {
	s := ToFirstUpper("abc")
	assert.Equal(t, "Abc", s)

	s = ToFirstUpper("ABC")
	assert.Equal(t, "Abc", s)

	s = ToFirstUpper("")
	assert.Equal(t, "", s)
}
