package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

// ForListingSAPExportFailures reads SAP export failure documents from storage.
type ForListingSAPExportFailures interface {
	List(limit int) ([]*models.SAPExportFailureLog, error)
}

// ForViewingSAPExportFailures is the application usecase for admins viewing recent SAP export failures.
type ForViewingSAPExportFailures interface {
	List(limit int) ([]*models.SAPExportFailureLog, error)
}
