package usecases

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var ErrTimesheetUserNotFound = errors.New("timesheet event: no matching user for employee email")
var ErrIncomeForTimesheetNotFoundForPeriod = errors.New("income_for_timesheet: no record for this user and period")

type syncIncomeForTimesheetUsecase struct {
	incomeRepo ForGettingIncomeForTimesheet
	userRepo   ForGettingTimesheetUser
}

func NewSyncIncomeForTimesheetUsecase(incomeRepo ForGettingIncomeForTimesheet, userRepo ForGettingTimesheetUser) ForSyncingIncomeForTimesheet {
	return &syncIncomeForTimesheetUsecase{incomeRepo, userRepo}
}

func (u *syncIncomeForTimesheetUsecase) SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error {
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

	existing, err := u.incomeRepo.GetByUserYearMonth(user.ID.Hex(), evt.Year, time.Month(evt.Month))
	switch {
	case errors.Is(err, ErrIncomeForTimesheetNotFoundForPeriod):
		income := models.CreatePayroll(*user, req, "")
		record := &models.IncomeForTimesheet{Income: *income, Sites: sites}
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
