package utils

import (
	"math"
	"strconv"
)

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func RealFloat(f float64) float64 {
	return math.Round(f*100) / 100
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
