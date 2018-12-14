package utils

import (
	"regexp"
	"testing"
)

func TestRandStringBytes(t *testing.T) {
	n := 62
	s := RandStringBytes(n)
	if len(s) != n {
		t.Errorf("Expected length %d but got %d", n, len(s))
	}

	reg := regexp.MustCompile("^[0-9a-zA-Z]*$")
	if !reg.MatchString(s) {
		t.Errorf("Expected string 0-9a-zA-Z but got %s", s)
	}
}
