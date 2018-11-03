package utils

import (
	"fmt"
	"math"
	"strconv"
)

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func FloatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func RealFloat(f float64) float64 {
	return math.Round(f*100) / 100
}
