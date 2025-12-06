package usecases

import (
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/entity"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestUsecaseExportIncome(t *testing.T) {
	t.Run("export individual income current month success", func(t *testing.T) {
		usecase, ctrl, mockRepoIncome := CreateExportIncomeUsecaseWithMock(t)
		defer ctrl.Finish()
		incomes := []*models.Income{
			&entity.MockIncome,
			&entity.MockIncome2,
		}
		mockRepoIncome.ExpectGetAllIncomeOfCurrentMonthByRole(incomes, time.Now())

		filename, err := usecase.ExportIncome("individual", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export individual income includes income from friendslog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		incomes := []*models.Income{
			&entity.MockIncome,
			&entity.MockIncome2,
			{ID: bson.ObjectIdHex("5bd1fda30fd2df2a3e41e571"), Role: "individual", WorkDate: "20", DailyRate: 750},
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.ExpectGetAllIncomeOfCurrentMonthByRole(incomes, time.Now())

		usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncome("individual", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export corporate income previous month success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		incomes := []*models.Income{
			&entity.MockIncome,
			&entity.MockIncome2,
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.ExpectGetAllIncomeOfPreviousMonthByRole(incomes)

		usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncome("corporate", "1")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}

func TestUsecaseExportIncomeSAPByStartDateAndEndDate(t *testing.T) {
	t.Run("export individual income SAP by start date and end date success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		incomes := []*models.Income{
			&entity.MockSoloCorporateIncome,
			&entity.MockSwardCorporateIncome,
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.GetAllIncomeByRoleStartDateAndEndDate(incomes, "individual", startDate, endDate)

		usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("individual", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export individual income SAP works with emoji", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		i := deepClone(&entity.MockSoloCorporateIncome)
		i.Name = "บจก. โซโล่ เลเวลลิ่ง 🦄"
		incomes := []*models.Income{
			i,
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.GetAllIncomeByRoleStartDateAndEndDate(incomes, "individual", startDate, endDate)

		usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("individual", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export corporate income SAP by start date and end date success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		incomes := []*models.Income{
			&entity.MockSoloCorporateIncome,
			&entity.MockSwardCorporateIncome,
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.GetAllIncomeByRoleStartDateAndEndDate(incomes, "corporate", startDate, endDate)

		usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("corporate", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}
