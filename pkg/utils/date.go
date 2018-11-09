package utils

import (
	"fmt"
	"time"
)

func GetNow() string {
	return time.Now().Format(time.RFC3339)
}

func GetCurrentMonth() string {
	y, m, _ := time.Now().Date()
	cm := fmt.Sprintf("%d-%d", y, int(m))
	return cm
}

func GetYearMonthNow() (int, time.Month) {
	t := time.Now()
	return t.Year(), t.Month()
}
