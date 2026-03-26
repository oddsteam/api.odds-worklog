package models

import "fmt"

// SAPExportRowError wraps an encoding/write failure for a single SAP line (TXN or WHT) so callers can log income identity.
type SAPExportRowError struct {
	RowIndex        int
	IncomeID        string
	UserID          string
	BankAccountName string
	LineKind        string
	Err             error
}

func (e *SAPExportRowError) Error() string {
	return fmt.Sprintf("sap export row %d (%s line) incomeId=%s userId=%s: %v",
		e.RowIndex, e.LineKind, e.IncomeID, e.UserID, e.Err)
}

func (e *SAPExportRowError) Unwrap() error {
	return e.Err
}
