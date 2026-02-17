package usecases

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestUsecaseExportIncome(t *testing.T) {
	t.Run("export individual income current month success", func(t *testing.T) {
		usecase, ctrl, mockRepoIncome := CreateExportIncomeUsecaseWithMock(t)
		defer ctrl.Finish()
		incomes := []*models.Income{
			&models.MockIncome,
			&models.MockIncome2,
		}
		mockRepoIncome.ExpectGetAllIncomeOfCurrentMonthByRole(incomes, time.Now())
		mockRepoIncome.ExpectGetStudentLoans()
		mockRepoIncome.ExpectAddExport()

		filename, err := usecase.ExportIncome("individual", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export corporate income previous month success", func(t *testing.T) {
		usecase, ctrl, mockRepoIncome := CreateExportIncomeUsecaseWithMock(t)
		defer ctrl.Finish()
		incomes := []*models.Income{
			&models.MockIncome,
			&models.MockIncome2,
		}
		mockRepoIncome.ExpectGetAllIncomeOfPreviousMonthByRole(incomes)
		mockRepoIncome.ExpectGetStudentLoans()
		mockRepoIncome.ExpectAddExport()

		filename, err := usecase.ExportIncome("corporate", "1")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}

func TestUsecaseExportIncomeSAPByStartDateAndEndDate(t *testing.T) {
	t.Run("export individual income SAP by start date and end date success", func(t *testing.T) {
		usecase, ctrl, mockRepoIncome := CreateExportIncomeUsecaseWithMock(t)
		defer ctrl.Finish()
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		incomes := []*models.Income{
			&models.MockSoloCorporateIncome,
			&models.MockSwardCorporateIncome,
		}
		mockRepoIncome.ExpectGetAllIncomeByRoleStartDateAndEndDate(incomes, "individual", startDate, endDate)
		mockRepoIncome.ExpectGetStudentLoans()
		mockRepoIncome.ExpectAddExport()

		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("individual", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export individual income SAP works with emoji", func(t *testing.T) {
		usecase, ctrl, mockRepoIncome := CreateExportIncomeUsecaseWithMock(t)
		defer ctrl.Finish()
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		i := deepClone(&models.MockSoloCorporateIncome)
		i.Name = "บจก. โซโล่ เลเวลลิ่ง 🦄"
		incomes := []*models.Income{
			i,
		}
		mockRepoIncome.ExpectGetAllIncomeByRoleStartDateAndEndDate(incomes, "individual", startDate, endDate)
		mockRepoIncome.ExpectGetStudentLoans()
		mockRepoIncome.ExpectAddExport()

		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("individual", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export corporate income SAP by start date and end date success", func(t *testing.T) {
		usecase, ctrl, mockRepoIncome := CreateExportIncomeUsecaseWithMock(t)
		defer ctrl.Finish()
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		incomes := []*models.Income{
			&models.MockSoloCorporateIncome,
			&models.MockSwardCorporateIncome,
		}
		mockRepoIncome.ExpectGetAllIncomeByRoleStartDateAndEndDate(incomes, "corporate", startDate, endDate)
		mockRepoIncome.ExpectGetStudentLoans()
		mockRepoIncome.ExpectAddExport()

		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("corporate", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}
