package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

func AddBlank(value string, length int) string {
	valLen := utf8.RuneCountInString(value)
	if valLen > length {
		return value
	}
	return value + strings.Repeat(" ", length-utf8.RuneCountInString(value))
}

func AmountStr(amt float64, n int) string {

	if amt == 0 && n <= 2 {
		return "0.00"
	}

	tempStr := fmt.Sprintf("%.2f", amt)

	l := len(tempStr)

	if l >= n {
		return tempStr[:n]
	}

	tStr := strings.Repeat("0", n-l)
	return tStr + tempStr
}

func LeftN(s string, n int) string {
	runes := []rune(s)
	if len(runes) >= n {
		return string(runes[:n])
	}
	return s
}

func ReceiveBRCode(bnkCode, brchCode string) string {
	var tempStr string

	switch bnkCode {
	case "002", "004", "006", "011", "014",
		"015", "020", "022", "024", "025",
		"033", "065", "066", "071", "073":
		if len(brchCode) >= 3 {
			tempStr = "0" + brchCode[:3]
		} else {
			tempStr = "0" + brchCode
		}

	case "030", "067", "072", "069", "052":
		if len(brchCode) >= 4 {
			tempStr = brchCode[:4]
		} else {
			tempStr = brchCode
		}

	case "034":
		tempStr = "0000"

	case "045":
		tempStr = "0010"

	default:
		tempStr = "0001"
	}

	return tempStr
}

func ReceiveAcCode(bnkCode, acCode string) (string, error) {
	acCode = strings.ReplaceAll(acCode, "-", "")
	acCode = strings.ReplaceAll(acCode, " ", "")
	var tempStr string
	if bnkCode == "011" {
		if len(acCode) == 10 {
			tempStr = "0" + acCode
		} else {
			tempStr = acCode
		}
		return tempStr, nil
	} else {
		return "", errors.New("not supported bank code")
	}

}

func FilterOthersThanThaiAndAscii(s string) string {
	runes := []rune(s)
	for j := range runes {
		r := runes[j]
		// Keep printable ASCII (32-126), and common Thai characters
		// Replace control characters and emojis with spaces
		if r < 32 || (r > 126 && r < 0x0E00) || (r > 0x0E7F && r < 0x2000) || r > 0x206F {
			runes[j] = 32
		}
	}
	return string(runes)
}
