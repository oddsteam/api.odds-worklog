package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentMonthInBuddistEra(t *testing.T) {
	n := time.Date(2022, time.Month(12), 1, 13, 30, 29, 0, time.UTC)
	actual := GetCurrentMonthInBuddistEra(n)
	assert.Equal(t, "12/2565", actual)
}

func TestGetCurrentMonthInBuddistEraWithOneDigitMonth(t *testing.T) {
	n := time.Date(2022, time.Month(1), 1, 13, 30, 29, 0, time.UTC)
	actual := GetCurrentMonthInBuddistEra(n)
	assert.Equal(t, "1/2565", actual)
}

func TestGetYearMonthInBuddistEra(t *testing.T) {
	n := time.Date(2022, time.Month(12), 1, 13, 30, 29, 0, time.UTC)
	y, m := GetYearMonthInBuddistEra(n)
	assert.Equal(t, 2565, y)
	assert.Equal(t, 12, m)
}
