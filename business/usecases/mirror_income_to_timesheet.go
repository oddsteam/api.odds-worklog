package usecases

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// mirrorIncomeToTimesheet upserts a manually entered Income into the income_from_timesheet
// collection for the given period, so that collection holds every income — not only the ones
// the timesheet consumer produced.
//
// The manual entry wins: every Income field is overwritten, including the workDate and
// workingHours the timesheet event owns. Only the per-site breakdown survives, because the
// income form carries no site data to replace it with; the next timesheet event for the period
// overwrites both again.
func mirrorIncomeToTimesheet(repo ForGettingIncomeFromTimesheet, income *models.Income, year int, month time.Month) error {
	existing, err := repo.GetByUserYearMonth(income.UserID, year, month)
	switch {
	case errors.Is(err, ErrIncomeFromTimesheetNotFoundForPeriod):
		return repo.Add(&models.IncomeFromTimesheet{Income: *income, Sites: []models.SiteWork{}})
	case err != nil:
		return err
	default:
		// Income carries the id, so copying it over would point the update at the income
		// collection's document instead of this one.
		id := existing.ID
		sites := existing.Sites
		existing.Income = *income
		existing.ID = id
		existing.Sites = sites
		return repo.Update(existing)
	}
}
