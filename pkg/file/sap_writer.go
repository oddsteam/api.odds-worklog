package file

import (
	"errors"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type sapWriter struct{}

func NewSAPWriter() *sapWriter { return &sapWriter{} }

func (w *sapWriter) WriteFile(name string, ics models.PayrollCycle, dateEff time.Time) (string, error) {
	strWrite, metas := ToSAP(ics, dateEff)

	if len(strWrite) == 0 {
		return "", errors.New("no data for export to SAP file")
	}

	file, filename, err := CreateFile(name)
	if err != nil {
		return "", err
	}
	encoder := charmap.Windows874.NewEncoder()
	writer := transform.NewWriter(file, encoder)
	defer file.Close()
	defer writer.Close()

	for i, record := range strWrite {
		incomeIdx := i / 2
		lineKind := "TXN"
		if i%2 == 1 {
			lineKind = "WHT"
		}
		m := metas[incomeIdx]
		row := createSAPRow(record)
		_, err := writer.Write([]byte(row))
		if err != nil {
			return "", &models.SAPExportRowError{
				RowIndex:        i,
				IncomeID:        m.IncomeID,
				UserID:          m.UserID,
				BankAccountName: m.BankAccountName,
				LineKind:        lineKind,
				Err:             err,
			}
		}
	}
	return filename, nil
}

func createSAPRow(record []string) string {
	return strings.Join(record, "") + "\n"
}

func ToSAP(ics models.PayrollCycle, dateEff time.Time) ([][]string, []models.IncomeRowMeta) {
	return ics.ProcessRecords(func(index int, i models.Payroll) [][]string {
		txn, wht := exportSAP(i, dateEff)
		return [][]string{txn, wht}
	})
}

func exportSAP(i models.Payroll, dateEff time.Time) ([]string, []string) {
	txn := toTransaction(i, dateEff)
	return txn.ToTXNLine(), txn.ToWHTLine()
}
