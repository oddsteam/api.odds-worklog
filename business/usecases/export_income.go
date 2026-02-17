package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type usecase struct {
	readRepo  ForGettingIncomeData
	writeRepo ForControllingIncomeData
	csvWriter ForWritingCSVFile
	sapWriter ForWritingSAPFile
}

func NewExportIncomeUsecase(r ForGettingIncomeData, ex ForControllingIncomeData,
	csvW ForWritingCSVFile, sapW ForWritingSAPFile) ForUsingExportIncome {
	return &usecase{
		readRepo:  r,
		writeRepo: ex,
		csvWriter: csvW,
		sapWriter: sapW,
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

	studentLoanList := u.readRepo.GetStudentLoans()

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

	studentLoanList := u.readRepo.GetStudentLoans()

	pc := models.NewPayrollCycle(incomes, studentLoanList)

	filename, err := u.sapWriter.WriteFile(role, *pc, dateEff)
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
