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

func TestNumberFormat(t *testing.T) {
	assert.Equal(t, "", FormatCommas(""))
	assert.Equal(t, "-1", FormatCommas("-1"))
	assert.Equal(t, "-12", FormatCommas("-12"))
	assert.Equal(t, "-123", FormatCommas("-123"))
	assert.Equal(t, "-1,234", FormatCommas("-1234"))
	assert.Equal(t, "-12,345", FormatCommas("-12345"))
	assert.Equal(t, "-123,456", FormatCommas("-123456"))
	assert.Equal(t, "-1,234,567", FormatCommas("-1234567"))
	assert.Equal(t, "-12,345,678", FormatCommas("-12345678"))
	assert.Equal(t, "-123,456,789", FormatCommas("-123456789"))
	assert.Equal(t, "-1,234,567,890", FormatCommas("-1234567890"))
	assert.Equal(t, "1,234,567,890", FormatCommas("1234567890"))
	assert.Equal(t, "123,456,789", FormatCommas("123456789"))
	assert.Equal(t, "12,345,678", FormatCommas("12345678"))
	assert.Equal(t, "1,234,567", FormatCommas("1234567"))
	assert.Equal(t, "123,456", FormatCommas("123456"))
	assert.Equal(t, "12,345", FormatCommas("12345"))
	assert.Equal(t, "1,234", FormatCommas("1234"))
	assert.Equal(t, "123", FormatCommas("123"))
	assert.Equal(t, "12", FormatCommas("12"))
	assert.Equal(t, "1", FormatCommas("1"))

	assert.Equal(t, "-1,234,567,890.00", FormatCommas("-1234567890.00"))
	assert.Equal(t, "1,234,567,890.99", FormatCommas("1234567890.99"))

	assert.Equal(t, "-1,234,567,890", FormatCommas("-1234567890."))
	assert.Equal(t, "1,234,567,890", FormatCommas("1234567890."))
}
