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
