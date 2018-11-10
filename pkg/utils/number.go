package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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

func FormatCommas(num string) string {
	if num == "" {
		return num
	}

	sp := strings.Split(num, ".")
	bd, ad := sp[0], ""
	if len(sp) == 2 {
		ad = sp[1]
	}

	re := regexp.MustCompile("(\\d+)(\\d{3})")
	ln := (len(bd) - 1) / 3
	for i := 0; i < ln; i++ {
		bd = re.ReplaceAllString(bd, "$1,$2")
	}

	if len(sp) == 2 && ad != "" {
		return bd + "." + ad
	}
	return bd
}
