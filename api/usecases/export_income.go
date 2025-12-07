package usecases

import (
	"encoding/csv"
	"errors"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/entity"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type ForUsingExportIncome interface {
	ExportIncome(role string, beforeMonth string) (string, error)
	ExportIncomeByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error)
	ExportIncomeSAPByStartDateAndEndDate(role string, startDate, endDate time.Time, dateEff time.Time) (string, error)
}

type usecase struct {
	repo     ForGettingIncomeData
	export   ForControllingIncomeData
	userRepo user.Repository
}

func NewExportIncomeUsecase(r ForGettingIncomeData, ex ForControllingIncomeData, ur user.Repository) ForUsingExportIncome {
	return &usecase{r, ex, ur}
}

func (u *usecase) ExportIncome(role string, beforeMonth string) (string, error) {
	return u.exportIncome(role, beforeMonth)
}

func (u *usecase) exportIncome(role string, beforeMonth string) (string, error) {
	var t time.Time
	if beforeMonth == "0" {
		t = time.Now()
	} else {
		t = time.Now().AddDate(0, -1, 0)
	}
	startDate, endDate := utils.GetStartDateAndEndDate(t)
	return u.ExportIncomeByStartDateAndEndDate(role, startDate, endDate)
}

func (u *usecase) ExportIncomeByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}

	incomes, err := u.repo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	ics := entity.NewIncomes(incomes, studentLoanList)
	strWrite, _ := ics.ToCSV()

	if len(strWrite) == 1 {
		return "", errors.New("no data for export to CSV file")
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.export.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (u *usecase) ExportIncomeSAPByStartDateAndEndDate(role string, startDate, endDate time.Time, dateEff time.Time) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	encoder := charmap.Windows874.NewEncoder()
	writer := transform.NewWriter(file, encoder)
	defer file.Close()
	defer writer.Close()

	if err != nil {
		return "", err
	}

	incomes, err := u.repo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	ics := entity.NewIncomes(incomes, studentLoanList)

	strWrite, _ := ics.ToSAP(dateEff)

	if len(strWrite) == 0 {
		return "", errors.New("no data for export to SAP file")
	}

	for _, record := range strWrite {
		row := createSAPRow(record)
		_, err := writer.Write([]byte(row))
		if err != nil {
			return "", err
		}
	}

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.export.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func createSAPRow(record []string) string {
	// r := filterEmojis(record)
	return strings.Join(record, "") + "\n"
}
