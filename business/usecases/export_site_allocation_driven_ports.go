package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForWritingSiteAllocationFile interface {
	WriteFile(name string, allocations []models.SiteAllocation) (string, error)
}
