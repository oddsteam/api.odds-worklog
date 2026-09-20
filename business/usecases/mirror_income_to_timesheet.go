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
//
// The period is stored on the record and SubmitDate is anchored to it, exactly as the timesheet
// sync does: these rows share the collection, so they have to share its key. Editing an income
// from an earlier period would otherwise mirror it into the month the edit happened in, because
// UpdatePayroll restamps SubmitDate with now.
func mirrorIncomeToTimesheet(repo ForGettingIncomeFromTimesheet, income *models.Income, year int, month time.Month) error {
	existing, err := repo.GetByUserYearMonth(income.UserID, year, month)
	switch {
	case errors.Is(err, ErrIncomeFromTimesheetNotFoundForPeriod):
		mirrored := *income
		mirrored.SubmitDate = timesheetPeriodSubmitDate(year, int(month))
		return repo.Add(&models.IncomeFromTimesheet{
			Income: mirrored,
			Year:   year,
			Month:  int(month),
			Sites:  []models.SiteWork{},
		})
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
		existing.SubmitDate = timesheetPeriodSubmitDate(year, int(month))
		existing.Year = year
		existing.Month = int(month)
		return repo.Update(existing)
	}
}
