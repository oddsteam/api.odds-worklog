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
	readStudentLoanRepo ForListStudentLoansInTheMonth
}

func NewExportIncomeUsecase(r ForGettingIncomeDataInTheMonth, ex ForLoggingExport, sapFail ForLoggingSAPExportFailure,
	csvW ForWritingCSVFile, sapW ForWritingSAPFile, rsl ForListStudentLoansInTheMonth) ForUsingExportIncome {
	return &usecase{
		readRepo:            r,
		writeRepo:           ex,
		sapExportFailureLog: sapFail,
		csvWriter:           csvW,
		sapWriter:           sapW,
		readStudentLoanRepo: rsl,
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
