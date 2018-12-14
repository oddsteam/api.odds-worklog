package utils

import (
	"fmt"
	"os"
	"time"
)

func CreateCVSFile(role string) (*os.File, string, error) {
	r := RandStringBytes(8)
	t := time.Now()
	tf := fmt.Sprintf("%d_%02d_%02d_%02d_%02d_%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	path := "files"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	filename := fmt.Sprintf("%s/%s_%s_%s.csv", path, role, tf, r)
	file, error := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, filename, error
}
