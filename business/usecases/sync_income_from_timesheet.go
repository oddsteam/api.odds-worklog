package usecases

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var ErrTimesheetUserNotFound = errors.New("timesheet event: no matching user for employee email")
var ErrIncomeFromTimesheetNotFoundForPeriod = errors.New("income_from_timesheet: no record for this user and period")

// hoursPerWorkDay converts the timesheet's day-based overtime figures into the hour-based
// units the worklog/payroll model expects, and derives the OT hourly rate from the user's
// daily rate the same way.
const hoursPerWorkDay = 8

type syncIncomeFromTimesheetUsecase struct {
	incomeRepo   ForGettingIncomeFromTimesheet
	userRepo     ForGettingTimesheetUser
	eventLogRepo ForLoggingTimesheetEvent
	siteRepo     ForGettingSiteByID
}

func NewSyncIncomeFromTimesheetUsecase(incomeRepo ForGettingIncomeFromTimesheet, userRepo ForGettingTimesheetUser, eventLogRepo ForLoggingTimesheetEvent, siteRepo ForGettingSiteByID) ForSyncingIncomeFromTimesheet {
	return &syncIncomeFromTimesheetUsecase{incomeRepo, userRepo, eventLogRepo, siteRepo}
}

func (u *syncIncomeFromTimesheetUsecase) SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error {
	if err := u.eventLogRepo.Save(evt); err != nil {
		return err
	}

	user, err := u.userRepo.GetByEmail(evt.Employee.Email)
	if err != nil {
		return err
	}
	attachSite(user, u.siteRepo)

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

	dailyRate, err := models.StringToFloat64(user.DailyIncome)
	if err != nil {
		dailyRate = 0
	}

	req := models.IncomeReq{
		WorkDate:      models.FloatToString(workingDays),
		WorkingHours:  models.FloatToString(overtimeDays * hoursPerWorkDay),
		SpecialIncome: models.FloatToString(dailyRate / hoursPerWorkDay),
	}

	existing, err := u.incomeRepo.GetByUserYearMonth(user.ID.Hex(), evt.Year, time.Month(evt.Month))
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
