package file

import (
	"encoding/csv"
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/entity"
)

type csvWriter struct{}

func NewCSVWriter() *csvWriter {
	return &csvWriter{}
}

func (w *csvWriter) WriteFile(name string, ics entity.Incomes) (string, error) {
	strWrite, _ := ics.ToCSV()

	if len(strWrite) == 1 {
		return "", errors.New("no data for export to CSV file")
	}

	file, filename, err := CreateFile(name)

	if err != nil {
		return "", err
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()
	defer file.Close()
	return filename, nil
}
