package utils

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNow(t *testing.T) {
	n := GetNow()
	r := regexp.MustCompile("\\d{4}-\\d{0,2}-\\d{0,2}T\\d{2}:\\d{2}:\\d{2}")
	assert.True(t, r.MatchString(n))
}

func TestGetNowAgain(t *testing.T) {
	n := time.Date(2022, time.Month(12), 1, 13, 30, 29, 0, time.UTC)
	actual := getNow(n)
	assert.Equal(t, "2022-12-01T13:30:29Z", actual)
}

func TestGetCurrentMonth(t *testing.T) {
	r := regexp.MustCompile("\\d{4}-\\d{0,2}")
	cm := GetCurrentMonth()
	assert.True(t, r.MatchString(cm))
}

func TestGetCurrentMonthAgain(t *testing.T) {
	n := time.Date(2022, time.Month(12), 1, 13, 30, 29, 0, time.UTC)
	actual := getCurrentMonth(n)
	assert.Equal(t, "2022-12", actual)
}

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

func TestGetYearMonthStringInBuddistEra(t *testing.T) {
	n := time.Date(2022, time.Month(1), 1, 13, 30, 29, 0, time.UTC)
	y, m := GetYearMonthStringInBuddistEra(n)
	assert.Equal(t, "2565", y)
	assert.Equal(t, "01", m)
}

func TestGetYearMonthNow(t *testing.T) {
	tn := time.Now()
	y, m := GetYearMonthNow()

	assert.Equal(t, tn.Year(), y)
	assert.Equal(t, tn.Month(), m)
}

func TestGetTheStartDateAndEndDateOfThisMonth(t *testing.T) {
	now := time.Date(2022, time.Month(1), 1, 13, 30, 29, 0, time.UTC)
	startDate, endDate := GetStartDateAndEndDate(now)

	expectedStartDate := time.Date(2022, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	expectedEndDate := time.Date(2022, time.Month(1), 31, 23, 59, 59, 0, time.UTC)
	assert.Equal(t, expectedStartDate, startDate)
	assert.Equal(t, expectedEndDate, endDate)
}
