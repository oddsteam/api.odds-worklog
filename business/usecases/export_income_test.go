package usecases

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
)

type stubSAPWriter struct {
	err error
}

func (s *stubSAPWriter) WriteFile(name string, ics models.PayrollCycle, dateEff time.Time) (string, error) {
	return "", s.err
}

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
		i.BankAccountName = "บจก. โซโล่ เลเวลลิ่ง 🦄"
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

	t.Run("export SAP persists failure log when writer returns SAPExportRowError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRead := mock_usecases.NewMockForGettingIncomeData(ctrl)
		mockWrite := mock_usecases.NewMockForControllingIncomeData(ctrl)
		mockSapFail := mock_usecases.NewMockForLoggingSAPExportFailure(ctrl)

		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		inc := models.MockSoloCorporateIncome
		incomes := []*models.Income{&inc}
		mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate("individual", startDate, endDate).Return(incomes, nil)
		mockRead.EXPECT().GetStudentLoans().Return(models.StudentLoanList{})

		underlying := errors.New("encoding: rune not supported by encoding")
		rowErr := &models.SAPExportRowError{
			RowIndex:        1,
			IncomeID:        inc.ID.Hex(),
			UserID:          inc.UserID,
			BankAccountName: inc.BankAccountName,
			LineKind:        "WHT",
			Err:             underlying,
		}
		mockSapFail.EXPECT().LogSAPExportFailure(gomock.Any()).Do(func(log *models.SAPExportFailureLog) {
			assert.Equal(t, inc.ID.Hex(), log.IncomeID)
			assert.Equal(t, inc.UserID, log.UserID)
			assert.Equal(t, "individual", log.Role)
			assert.Equal(t, 1, log.RowIndex)
			assert.Equal(t, "WHT", log.LineKind)
			assert.Contains(t, log.ErrorMessage, "sap export row")
		}).Return(nil)

		u := NewExportIncomeUsecase(mockRead, mockWrite, mockSapFail, file.NewCSVWriter(), &stubSAPWriter{err: rowErr}, mockRead)
		_, err := u.ExportIncomeSAPByStartDateAndEndDate("individual", startDate, endDate, dateEff)
		assert.Error(t, err)
		assert.ErrorIs(t, err, underlying)
	})
}
