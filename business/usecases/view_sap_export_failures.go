package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

const (
	defaultSapExportFailureLimit = 100
	maxSapExportFailureLimit     = 200
)

type viewSAPExportFailuresUsecase struct {
	list ForListingSAPExportFailures
}

func NewViewSAPExportFailuresUsecase(list ForListingSAPExportFailures) ForViewingSAPExportFailures {
	return &viewSAPExportFailuresUsecase{list: list}
}

func normalizeSapExportFailureLimit(limit int) int {
	if limit <= 0 {
		return defaultSapExportFailureLimit
	}
	if limit > maxSapExportFailureLimit {
		return maxSapExportFailureLimit
	}
	return limit
}

func (u *viewSAPExportFailuresUsecase) List(limit int) ([]*models.SAPExportFailureLog, error) {
	n := normalizeSapExportFailureLimit(limit)
	return u.list.List(n)
}
