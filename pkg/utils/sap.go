package utils

import (
	"errors"
	"fmt"
	"math"
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
	// 1. ปัดเศษให้เหลือ 2 ตำแหน่งทศนิยม (Round)
	// Go ไม่มีฟังก์ชัน Round() ในตัวสำหรับ float64, จึงต้องสร้างเอง
	// Note: การใช้ math.Round(x * 100) / 100 เป็นวิธีปัดเศษมาตรฐานทางคณิตศาสตร์
	roundedAmt := math.Ceil(amt*100) / 100
	fmt.Printf("%.2f", roundedAmt)
	// 2. Format ตัวเลขให้อยู่ในรูปแบบ "XX.XX"
	// fmt.Sprintf("%.2f", ...) จะจัดการเรื่องการมีทศนิยม 2 ตำแหน่งให้
	tempStr := fmt.Sprintf("%.2f", roundedAmt)

	// 3. คำนวณความยาวปัจจุบันของสตริง
	l := len(tempStr)

	// 4. เติม '0' นำหน้าจนครบความยาว n
	if l >= n {
		// ถ้าความยาวเกินหรือเท่ากับ n ก็ตัดส่วนเกินออก (ถ้าจำเป็น)
		return tempStr[:n]
	}

	// จำนวนศูนย์ที่ต้องเติม = n - l
	tStr := strings.Repeat("0", n-l)

	// 5. รวมสตริงที่เติมศูนย์นำหน้ากับสตริงตัวเลข
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
