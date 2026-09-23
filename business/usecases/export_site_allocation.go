package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// exportSiteAllocationUsecase reports how the month's income splits across client sites.
// It reads income_from_timesheet directly rather than through the shared income source the
// other exports use, because that source unwraps the records to plain Income and the
// per-site day counts this report is built from live outside Income.
type exportSiteAllocationUsecase struct {
	readRepo  ForGettingIncomeFromTimesheetInTheMonth
	exportLog ForLoggingExport
	writer    ForWritingSiteAllocationFile
}

func NewExportSiteAllocationUsecase(r ForGettingIncomeFromTimesheetInTheMonth, ex ForLoggingExport, w ForWritingSiteAllocationFile) ForUsingExportSiteAllocation {
	return &exportSiteAllocationUsecase{readRepo: r, exportLog: ex, writer: w}
}

func (u *exportSiteAllocationUsecase) ExportSiteAllocation(role string, monthIndex string) (string, error) {
	t := time.Now()
	if monthIndex != "0" {
		t = t.AddDate(0, -1, 0)
	}
	startDate, endDate := models.GetStartDateAndEndDate(t)
	return u.ExportSiteAllocationByStartDateAndEndDate(role, startDate, endDate)
}

func (u *exportSiteAllocationUsecase) ExportSiteAllocationByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error) {
	records, err := u.readRepo.GetAllByRoleStartDateAndEndDate(role, startDate, endDate)
	if err != nil {
		return "", err
	}

	filename, err := u.writer.WriteFile("site_allocation_"+role, models.BuildSiteAllocations(records))
	if err != nil {
		return "", err
	}

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	if err := u.exportLog.AddExport(&ep); err != nil {
		return "", err
	}

	return filename, nil
}
