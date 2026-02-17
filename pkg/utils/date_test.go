package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetYearMonthStringInBuddistEra(t *testing.T) {
	n := time.Date(2022, time.Month(1), 1, 13, 30, 29, 0, time.UTC)
	y, m := GetYearMonthStringInBuddistEra(n)
	assert.Equal(t, "2565", y)
	assert.Equal(t, "01", m)
}
