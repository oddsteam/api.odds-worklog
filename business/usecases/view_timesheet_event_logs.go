package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

const (
	defaultTimesheetEventLogLimit = 100
	maxTimesheetEventLogLimit     = 200
)

type viewTimesheetEventLogsUsecase struct {
	list ForListingTimesheetEventLogs
}

func NewViewTimesheetEventLogsUsecase(list ForListingTimesheetEventLogs) ForListingTimesheetEventLogs {
	return &viewTimesheetEventLogsUsecase{list: list}
}

func normalizeTimesheetEventLogLimit(limit int) int {
	if limit <= 0 {
		return defaultTimesheetEventLogLimit
	}
	if limit > maxTimesheetEventLogLimit {
		return maxTimesheetEventLogLimit
	}
	return limit
}

func (u *viewTimesheetEventLogsUsecase) List(limit int) ([]*models.TimesheetEventLog, error) {
	n := normalizeTimesheetEventLogLimit(limit)
	return u.list.List(n)
}
