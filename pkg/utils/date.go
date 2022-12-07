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

func GetCurrentMonth() string {
	return getCurrentMonth(time.Now())
}

func getCurrentMonth(now time.Time) string {
	y, m, _ := now.Date()
	cm := fmt.Sprintf("%d-%d", y, int(m))
	return cm
}

func GetCurrentMonthInBuddistEra(now time.Time) string {
	y, m, _ := now.Date()
	cm := fmt.Sprintf("%d/%d", int(m), y+543)
	return cm
}

func GetYearMonthNow() (int, time.Month) {
	t := time.Now()
	return t.Year(), t.Month()
}
