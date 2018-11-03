package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToFloat(t *testing.T) {
	f, err := StringToFloat64("100000")
	assert.NoError(t, err)
	assert.Equal(t, 100000.0, f)

	f, err = StringToFloat64("1234.567890")
	assert.NoError(t, err)
	assert.Equal(t, 1234.567890, f)

	f, err = StringToFloat64("1234.56789")
	assert.NoError(t, err)
	assert.Equal(t, 1234.56789, f)
}

func TestFloatToString(t *testing.T) {
	f := FloatToString(100000.0)
	assert.Equal(t, "100000.00", f)

	f = FloatToString(1234.565890)
	assert.Equal(t, "1234.57", f)

	f = FloatToString(1234.564)
	assert.Equal(t, "1234.56", f)
}

func TestRealFloat(t *testing.T) {
	f := RealFloat(100000.0)
	assert.Equal(t, 100000.00, f)

	f = RealFloat(1234.565890)
	assert.Equal(t, 1234.57, f)

	f = RealFloat(1234.564)
	assert.Equal(t, 1234.56, f)
}
