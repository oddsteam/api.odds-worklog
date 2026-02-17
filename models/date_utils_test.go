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

func TestGetYearMonthNow(t *testing.T) {
	tn := time.Now()
	y, m := GetYearMonthNow()

	assert.Equal(t, tn.Year(), y)
	assert.Equal(t, tn.Month(), m)
}

func TestGetStartDateAndEndDate(t *testing.T) {
	now := time.Date(2022, time.Month(1), 15, 13, 30, 29, 0, time.UTC)
	startDate, endDate := GetStartDateAndEndDate(now)

	expectedStartDate := time.Date(2022, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	expectedEndDate := time.Date(2022, time.Month(1), 31, 23, 59, 59, 0, time.UTC)
	assert.Equal(t, expectedStartDate, startDate)
	assert.Equal(t, expectedEndDate, endDate)
}
