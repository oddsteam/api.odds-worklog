package models

import (
	"fmt"
	"time"
)

func GetCurrentMonthInBuddistEra(now time.Time) string {
	y, m := GetYearMonthInBuddistEra(now)
	cm := fmt.Sprintf("%d/%d", m, y)
	return cm
}

func GetYearMonthInBuddistEra(now time.Time) (int, int) {
	y, m, _ := now.Date()
	return y + 543, int(m)
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

func GetYearMonthStringInBuddistEra(now time.Time) (string, string) {
	y, m, _ := now.Date()
	return fmt.Sprintf("%d", y+543), fmt.Sprintf("%02d", m)
}
