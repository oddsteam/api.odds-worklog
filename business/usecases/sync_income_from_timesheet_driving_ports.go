package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForSyncingIncomeFromTimesheet interface {
	SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error
}
