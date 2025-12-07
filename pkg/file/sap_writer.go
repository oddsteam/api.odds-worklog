package file

import (
	"errors"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/entity"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type sapWriter struct{}

func NewSAPWriter() *sapWriter { return &sapWriter{} }

func (w *sapWriter) WriteFile(name string, ics entity.Incomes, dateEff time.Time) (string, error) {
	strWrite, _ := ToSAP(ics, dateEff)

	if len(strWrite) == 0 {
		return "", errors.New("no data for export to SAP file")
	}

	file, filename, err := CreateFile(name)
	encoder := charmap.Windows874.NewEncoder()
	writer := transform.NewWriter(file, encoder)
	defer file.Close()
	defer writer.Close()

	if err != nil {
		return "", err
	}

	for _, record := range strWrite {
		row := createSAPRow(record)
		_, err := writer.Write([]byte(row))
		if err != nil {
			return "", err
		}
	}
	return filename, nil
}

func createSAPRow(record []string) string {
	return strings.Join(record, "") + "\n"
}

func ToSAP(ics entity.Incomes, dateEff time.Time) ([][]string, []string) {
	return ics.ProcessRecords(func(index int, i *entity.Income) [][]string {
		txn, wht := exportSAP(*i, dateEff)
		return [][]string{txn, wht}
	})
}

func exportSAP(i entity.Income, dateEff time.Time) ([]string, []string) {
	txn := toTransaction(i, dateEff)
	return txn.ToTXNLine(), txn.ToWHTLine()
}
