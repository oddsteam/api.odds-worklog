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

func TestParseDateSuccess(t *testing.T) {
	actual, _ := ParseDate("2022-12-01T13:30:29Z")
	expected := time.Date(2022, time.Month(12), 1, 13, 30, 29, 0, time.UTC)
	assert.Equal(t, expected, actual)
}
func TestParseAnotherDate(t *testing.T) {
	actual, _ := ParseDate("2024-07-26T06:26:25.531Z")
	expected := time.Date(2024, time.Month(7), 26, 6, 26, 25, 531000000, time.UTC)
	assert.Equal(t, expected, actual)
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

func TestGetYearMonthStringInBuddistEra(t *testing.T) {
	n := time.Date(2022, time.Month(1), 1, 13, 30, 29, 0, time.UTC)
	y, m := GetYearMonthStringInBuddistEra(n)
	assert.Equal(t, "2565", y)
	assert.Equal(t, "01", m)
}
