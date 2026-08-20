package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type exportIncomeFromTimesheetUsecase struct {
	incomeRepo      ForGettingIncomeFromTimesheetInTheMonth
	studentLoanRepo ForListStudentLoansInTheMonth
	userRepo        ForListingUsersByRole
	siteRepo        ForListingSites
	csvWriter       ForWritingCSVFile
}

func NewExportIncomeFromTimesheetUsecase(incomeRepo ForGettingIncomeFromTimesheetInTheMonth, studentLoanRepo ForListStudentLoansInTheMonth, userRepo ForListingUsersByRole, siteRepo ForListingSites, csvWriter ForWritingCSVFile) ForUsingExportIncomeFromTimesheet {
	return &exportIncomeFromTimesheetUsecase{incomeRepo, studentLoanRepo, userRepo, siteRepo, csvWriter}
}

func (u *exportIncomeFromTimesheetUsecase) ExportIncomeFromTimesheet(role string, monthIndex string) (string, error) {
	var t time.Time
	if monthIndex == "0" {
		t = time.Now()
	} else {
		t = time.Now().AddDate(0, -1, 0)
	}
	startDate, endDate := models.GetStartDateAndEndDate(t)

	return u.ExportIncomeFromTimesheetByStartDateAndEndDate(role, startDate, endDate)
}

func (u *exportIncomeFromTimesheetUsecase) ExportIncomeFromTimesheetByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error) {
	records, err := u.incomeRepo.GetAllByRoleStartDateAndEndDate(role, startDate, endDate)
	if err != nil {
		return "", err
	}

	incomes := make([]*models.Income, len(records))
	for i, record := range records {
		incomes[i] = &record.Income
	}

	u.enrichSiteNames(role, incomes)

	studentLoanList := u.studentLoanRepo.GetStudentLoans()

	pc := models.NewPayrollCycle(incomes, studentLoanList)
	return u.csvWriter.WriteFile("individual-timesheet", *pc)
}

func (u *exportIncomeFromTimesheetUsecase) enrichSiteNames(role string, incomes []*models.Income) {
	if len(incomes) == 0 {
		return
	}

	users, err := u.userRepo.GetByRole(role)
	if err != nil || len(users) == 0 {
		return
	}

	sites, err := u.siteRepo.GetSiteGroup()
	if err != nil {
		return
	}

	siteNameByID := make(map[string]string, len(sites))
	for _, site := range sites {
		if site == nil {
			continue
		}
		siteNameByID[site.ID.Hex()] = site.Name
	}

	siteNameByUserID := make(map[string]string, len(users))
	for _, user := range users {
		if user == nil || user.SiteID == "" {
			continue
		}
		if name, ok := siteNameByID[user.SiteID]; ok {
			siteNameByUserID[user.ID.Hex()] = name
		}
	}

	for _, income := range incomes {
		if income == nil {
			continue
		}
		income.SiteName = siteNameByUserID[income.UserID]
	}
}
