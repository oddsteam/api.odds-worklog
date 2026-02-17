package usecases

import "time"

type ForUsingExportIncome interface {
	ExportIncome(role string, monthIndex string) (string, error)
	ExportIncomeByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error)
	ExportIncomeSAPByStartDateAndEndDate(role string, startDate, endDate time.Time, dateEff time.Time) (string, error)
}
