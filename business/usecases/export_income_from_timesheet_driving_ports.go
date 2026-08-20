package usecases

import "time"

type ForUsingExportIncomeFromTimesheet interface {
	ExportIncomeFromTimesheet(role string, monthIndex string) (string, error)
	ExportIncomeFromTimesheetByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error)
}
