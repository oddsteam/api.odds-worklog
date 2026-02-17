package utils

import (
	"fmt"
	"time"
)

func GetYearMonthStringInBuddistEra(now time.Time) (string, string) {
	y, m, _ := now.Date()
	return fmt.Sprintf("%d", y+543), fmt.Sprintf("%02d", m)
}
