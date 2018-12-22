package utils

import "strings"

func ToFirstUpper(s string) string {
	return strings.Title(strings.ToLower(s))
}
