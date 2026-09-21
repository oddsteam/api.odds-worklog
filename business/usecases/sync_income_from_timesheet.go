package usecases

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var ErrTimesheetUserNotFound = errors.New("timesheet event: no matching user for employee email")
var ErrTimesheetEventOutOfPeriod = errors.New("timesheet event: not for the current month")
var ErrIncomeFromTimesheetNotFoundForPeriod = errors.New("income_from_timesheet: no record for this user and period")

type syncIncomeFromTimesheetUsecase struct {
	incomeRepo   ForGettingIncomeFromTimesheet
	userRepo     ForGettingTimesheetUser
	eventLogRepo ForLoggingTimesheetEvent
}

func NewSyncIncomeFromTimesheetUsecase(incomeRepo ForGettingIncomeFromTimesheet, userRepo ForGettingTimesheetUser, eventLogRepo ForLoggingTimesheetEvent) ForSyncingIncomeFromTimesheet {
	return &syncIncomeFromTimesheetUsecase{incomeRepo, userRepo, eventLogRepo}
}

func (u *syncIncomeFromTimesheetUsecase) SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error {
	if err := u.eventLogRepo.Save(evt); err != nil {
		return err
	}

	// A record is stamped with submitDate at save time (same as the manual income flow),
	// never with the event's year/month — so an event for an older month has nowhere of
	// its own to land and would only overwrite the current month's record. Drop it: the
	// raw payload above is still kept for audit.
	year, month := models.GetYearMonthNow()
	if evt.Year != year || time.Month(evt.Month) != month {
		return ErrTimesheetEventOutOfPeriod
	}

	user, err := u.userRepo.GetByEmail(evt.Employee.Email)
	if err != nil {
		return err
	}

	var workingDays, overtimeDays float64
	sites := make([]models.SiteWork, 0, len(evt.Sites))
	for _, s := range evt.Sites {
		workingDays += s.WorkingDays
		overtimeDays += s.OvertimeDays
		sites = append(sites, models.SiteWork{
			ClientSite:   s.ClientSite,
			CustomerName: s.CustomerName,
			WorkingDays:  s.WorkingDays,
			OvertimeDays: s.OvertimeDays,
		})
	}

	req := models.IncomeReq{
		WorkDate:      models.FloatToString(workingDays),
		WorkingHours:  models.FloatToString(overtimeDays),
		SpecialIncome: "0",
	}

	existing, err := u.incomeRepo.GetByUserYearMonth(user.ID.Hex(), year, month)
	switch {
	case errors.Is(err, ErrIncomeFromTimesheetNotFoundForPeriod):
		income := models.CreatePayroll(*user, req, "")
		record := &models.IncomeFromTimesheet{Income: *income, Sites: sites}
		if err := u.incomeRepo.Add(record); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		models.UpdatePayroll(*user, req, existing.Note, &existing.Income)
		existing.Sites = sites
		if err := u.incomeRepo.Update(existing); err != nil {
			return err
		}
	}

	return nil
}
