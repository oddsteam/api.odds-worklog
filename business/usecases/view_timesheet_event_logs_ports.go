package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

// ForListingTimesheetEventLogs reads timesheet event log documents. It is both the storage port the
// usecase depends on and the port the HTTP handler is given — the usecase decorates a repository
// with limit normalization, so both sides speak the same contract.
type ForListingTimesheetEventLogs interface {
	List(limit int) ([]*models.TimesheetEventLog, error)
}
