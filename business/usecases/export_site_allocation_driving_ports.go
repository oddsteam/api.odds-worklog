package usecases

import "time"

type ForUsingExportSiteAllocation interface {
	ExportSiteAllocation(role string, monthIndex string) (string, error)
	ExportSiteAllocationByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error)
}
