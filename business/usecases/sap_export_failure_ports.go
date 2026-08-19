package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

// ForListingSAPExportFailures reads SAP export failure documents. It is both the storage port the
// usecase depends on and the port the HTTP handler is given — the usecase decorates a repository
// with limit normalization, so both sides speak the same contract.
type ForListingSAPExportFailures interface {
	List(limit int) ([]*models.SAPExportFailureLog, error)
}
