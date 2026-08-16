package usecases

import (
	"errors"
	"log"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type usecase struct {
	readRepo            ForGettingIncomeDataInTheMonth
	writeRepo           ForLoggingExport
	sapExportFailureLog ForLoggingSAPExportFailure
	csvWriter           ForWritingCSVFile
	sapWriter           ForWritingSAPFile
	peakWriter          ForWritingPeakFile
	readStudentLoanRepo ForListStudentLoansInTheMonth
	userRepo            ForListingUsersByRole
	siteRepo            ForListingSites
}

func NewExportIncomeUsecase(r ForGettingIncomeDataInTheMonth, ex ForLoggingExport, sapFail ForLoggingSAPExportFailure,
	csvW ForWritingCSVFile, sapW ForWritingSAPFile, peakW ForWritingPeakFile, rsl ForListStudentLoansInTheMonth,
	userRepo ForListingUsersByRole, siteRepo ForListingSites) ForUsingExportIncome {
	return &usecase{
		readRepo:            r,
		writeRepo:           ex,
		sapExportFailureLog: sapFail,
		csvWriter:           csvW,
		sapWriter:           sapW,
		peakWriter:          peakW,
		readStudentLoanRepo: rsl,
		userRepo:            userRepo,
		siteRepo:            siteRepo,
	}
}

func (u *usecase) ExportIncome(role string, monthIndex string) (string, error) {
	var t time.Time
	if monthIndex == "0" {
		t = time.Now()
	} else {
		t = time.Now().AddDate(0, -1, 0)
	}
	startDate, endDate := models.GetStartDateAndEndDate(t)
	return u.ExportIncomeByStartDateAndEndDate(role, startDate, endDate)
}

func (u *usecase) ExportIncomeByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error) {
	incomes, err := u.readRepo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	u.enrichSiteNames(role, incomes)

	studentLoanList := u.readStudentLoanRepo.GetStudentLoans()

	pc := models.NewPayrollCycle(incomes, studentLoanList)
	filename, err := u.csvWriter.WriteFile(role, *pc)
	if err != nil {
		return "", err
	}

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.writeRepo.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (u *usecase) ExportIncomeSAPByStartDateAndEndDate(role string, startDate, endDate time.Time, dateEff time.Time) (string, error) {
	incomes, err := u.readRepo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	studentLoanList := u.readStudentLoanRepo.GetStudentLoans()

	pc := models.NewPayrollCycle(incomes, studentLoanList)

	filename, err := u.sapWriter.WriteFile(role, *pc, dateEff)
	if err != nil {
		var rowErr *models.SAPExportRowError
		if errors.As(err, &rowErr) {
			underlying := ""
			if rowErr.Err != nil {
				underlying = rowErr.Err.Error()
			}
			logEntry := &models.SAPExportFailureLog{
				CreatedAt:       time.Now(),
				Role:            role,
				StartDate:       startDate,
				EndDate:         endDate,
				DateEffective:   dateEff,
				IncomeID:        rowErr.IncomeID,
				UserID:          rowErr.UserID,
				BankAccountName: rowErr.BankAccountName,
				RowIndex:        rowErr.RowIndex,
				LineKind:        rowErr.LineKind,
				ErrorMessage:    rowErr.Error(),
				UnderlyingError: underlying,
			}
			if logErr := u.sapExportFailureLog.LogSAPExportFailure(logEntry); logErr != nil {
				log.Printf("sap export failure log: %v", logErr)
			}
		}
		return "", err
	}

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.writeRepo.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (u *usecase) ExportPeak(role string, monthIndex string) (string, error) {
	var t time.Time
	if monthIndex == "0" {
		t = time.Now()
	} else {
		t = time.Now().AddDate(0, -1, 0)
	}
	startDate, endDate := models.GetStartDateAndEndDate(t)
	return u.ExportPeakByStartDateAndEndDate(role, startDate, endDate)
}

func (u *usecase) ExportPeakByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error) {
	incomes, err := u.readRepo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)
	if err != nil {
		return "", err
	}

	filename, err := u.peakWriter.WriteFile("peak_"+role, incomes, u.peakCodesByUserID(role), startDate, time.Now())
	if err != nil {
		return "", err
	}

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.writeRepo.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (u *usecase) peakCodesByUserID(role string) map[string]string {
	codes := map[string]string{}
	if u.userRepo == nil {
		return codes
	}
	users, err := u.userRepo.GetByRole(role)
	if err != nil || len(users) == 0 {
		return codes
	}
	for _, user := range users {
		if user == nil {
			continue
		}
		codes[user.ID.Hex()] = user.PeakCode
	}
	return codes
}

func (u *usecase) enrichSiteNames(role string, incomes []*models.Income) {
	if u.userRepo == nil || u.siteRepo == nil || len(incomes) == 0 {
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
