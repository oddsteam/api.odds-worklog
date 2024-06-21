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
	y, m := GetYearMonthInBuddistEra(now)
	cm := fmt.Sprintf("%d/%d", m, y)
	return cm
}

func GetYearMonthInBuddistEra(now time.Time) (int, int) {
	y, m, _ := now.Date()
	return y + 543, int(m)
}

func GetYearMonthStringInBuddistEra(now time.Time) (string, string) {
	y, m, _ := now.Date()
	return fmt.Sprintf("%d", y+543), fmt.Sprintf("%02d", m)
}

func GetYearMonthNow() (int, time.Month) {
	t := time.Now()
	return t.Year(), t.Month()
}

func GetStartDateAndEndDate(now time.Time) (time.Time, time.Time) {
	return getStartDate(now), getEndDate(now)
}

func getStartDate(now time.Time) time.Time {
	beginningOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return beginningOfMonth
}

func getEndDate(now time.Time) time.Time {
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	endOfMonth := firstOfNextMonth.Add(-time.Second)
	return endOfMonth
}
