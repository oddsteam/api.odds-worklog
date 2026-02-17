package utils

import (
	"fmt"
	"time"
)

func GetNow() string {
	return getNow(time.Now())
}

func getNow(now time.Time) string {
	return now.Format(time.RFC3339)
}

func ParseDate(s string) (time.Time, error) {
	layout := time.RFC3339Nano
	return time.Parse(layout, s)
}

func GetCurrentMonth() string {
	return getCurrentMonth(time.Now())
}

func getCurrentMonth(now time.Time) string {
	y, m, _ := now.Date()
	cm := fmt.Sprintf("%d-%d", y, int(m))
	return cm
}

func GetYearMonthStringInBuddistEra(now time.Time) (string, string) {
	y, m, _ := now.Date()
	return fmt.Sprintf("%d", y+543), fmt.Sprintf("%02d", m)
}
