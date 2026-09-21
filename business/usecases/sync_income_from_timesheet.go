package usecases

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var ErrTimesheetUserNotFound = errors.New("timesheet event: no matching user for employee email")
var ErrTimesheetEventOutOfPeriod = errors.New("timesheet event: not for the current month")
var ErrIncomeFromTimesheetNotFoundForPeriod = errors.New("income_from_timesheet: no record for this user and period")

// hoursPerWorkDay converts the timesheet's day-based overtime figures into the hour-based
// units the worklog/payroll model expects, and derives the OT hourly rate from the user's
// daily rate the same way.
const hoursPerWorkDay = 8

// timesheetSpecialIncomeLineKind marks the error-log entries this usecase writes, so they can be
// told apart from the SAP export rows that share the collection.
const timesheetSpecialIncomeLineKind = "timesheet-special-income"

type syncIncomeFromTimesheetUsecase struct {
	incomeRepo     ForGettingIncomeFromTimesheet
	userRepo       ForGettingTimesheetUser
	eventLogRepo   ForLoggingTimesheetEvent
	siteRepo       ForGettingSiteByID
	failureLogRepo ForLoggingSAPExportFailure
}

func NewSyncIncomeFromTimesheetUsecase(incomeRepo ForGettingIncomeFromTimesheet, userRepo ForGettingTimesheetUser, eventLogRepo ForLoggingTimesheetEvent, siteRepo ForGettingSiteByID, failureLogRepo ForLoggingSAPExportFailure) ForSyncingIncomeFromTimesheet {
	return &syncIncomeFromTimesheetUsecase{incomeRepo, userRepo, eventLogRepo, siteRepo, failureLogRepo}
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

	overtimeHours := overtimeDays * hoursPerWorkDay

	// A user with OT but no usable daily rate can't have their special income calculated. Store
	// the rest of the record anyway (work days would otherwise be lost too) and surface the case
	// on the error log page so someone can fill the rate in and let the next event re-sync it.
	dailyRate, parseErr := models.StringToFloat64(user.DailyIncome)
	if parseErr != nil {
		dailyRate = 0
	}
	if overtimeHours > 0 && dailyRate <= 0 {
		u.logMissingDailyRate(user, evt, overtimeHours, parseErr)
	}

	req := models.IncomeReq{
		WorkDate:      models.FloatToString(workingDays),
		WorkingHours:  models.FloatToString(overtimeHours),
		SpecialIncome: models.FloatToString(dailyRate / hoursPerWorkDay),
	}

	existing, err := u.incomeRepo.GetByUserYearMonth(user.ID.Hex(), year, month)
	switch {
	case errors.Is(err, ErrIncomeFromTimesheetNotFoundForPeriod):
		income := models.CreatePayroll(*user, req, "")
		income.SiteName = "Timesheet"
		record := &models.IncomeFromTimesheet{Income: *income, Sites: sites}
		if err := u.incomeRepo.Add(record); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		models.UpdatePayroll(*user, req, existing.Note, &existing.Income)
		existing.SiteName = "Timesheet"
		existing.Sites = sites
		if err := u.incomeRepo.Update(existing); err != nil {
			return err
		}
	}

	return nil
}

func (u *syncIncomeFromTimesheetUsecase) logMissingDailyRate(user *models.User, evt models.TimesheetMonthlySummaryEvent, overtimeHours float64, cause error) {
	if u.failureLogRepo == nil {
		return
	}

	periodStart := time.Date(evt.Year, time.Month(evt.Month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, -1)
	underlying := ""
	if cause != nil {
		underlying = cause.Error()
	}

	entry := &models.SAPExportFailureLog{
		CreatedAt:       time.Now(),
		Role:            user.Role,
		StartDate:       periodStart,
		EndDate:         periodEnd,
		UserID:          user.ID.Hex(),
		BankAccountName: user.BankAccountName,
		LineKind:        timesheetSpecialIncomeLineKind,
		ErrorMessage: fmt.Sprintf(
			"timesheet sync %04d-%02d: %s has %s OT hours but no usable daily rate (dailyIncome=%q) — special income was stored as 0",
			evt.Year, evt.Month, user.Email, models.FloatToString(overtimeHours), user.DailyIncome,
		),
		UnderlyingError: underlying,
	}
	if err := u.failureLogRepo.LogSAPExportFailure(entry); err != nil {
		log.Printf("timesheet sync: missing daily rate log: %v", err)
	}
}
