package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealFloat(t *testing.T) {
	f := RealFloat(100000.0)
	assert.Equal(t, 100000.00, f)

	f = RealFloat(1234.565890)
	assert.Equal(t, 1234.57, f)

	f = RealFloat(1234.564)
	assert.Equal(t, 1234.56, f)
}

func TestIsNumber(t *testing.T) {
	assert.True(t, IsNumeric("1234567890"))
	assert.False(t, IsNumeric("1,234,567,890"))
	assert.False(t, IsNumeric(""))
}
