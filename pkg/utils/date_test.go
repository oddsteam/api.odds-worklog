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

func TestGetCurrentMonth(t *testing.T) {
	r := regexp.MustCompile("\\d{4}-\\d{0,2}")
	cm := GetCurrentMonth()
	assert.True(t, r.MatchString(cm))
}

func TestGetYearMonthNow(t *testing.T) {
	tn := time.Now()
	y, m := GetYearMonthNow()

	assert.Equal(t, tn.Year(), y)
	assert.Equal(t, tn.Month(), m)
}
