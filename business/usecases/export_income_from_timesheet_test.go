package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
)

type stubTimesheetExportUsersByRole struct{}

func (s *stubTimesheetExportUsersByRole) GetByRole(role string) ([]*models.User, error) {
	return []*models.User{}, nil
}

type stubTimesheetExportSites struct{}

func (s *stubTimesheetExportSites) GetSiteGroup() ([]*models.Site, error) {
	return []*models.Site{}, nil
}

type stubTimesheetExportStudentLoans struct{}

func (s *stubTimesheetExportStudentLoans) GetStudentLoans() models.StudentLoanList {
	return models.StudentLoanList{}
}

func TestExportIncomeFromTimesheet(t *testing.T) {
	t.Run("converts records to Income and writes a CSV via the existing PayrollCycle machinery", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		record := &models.IncomeFromTimesheet{Income: models.MockIncome}
		record.UserID = bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539").Hex()

		incomeRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", gomock.Any(), gomock.Any()).
			Return([]*models.IncomeFromTimesheet{record}, nil)

		uc := NewExportIncomeFromTimesheetUsecase(incomeRepo, &stubTimesheetExportStudentLoans{}, &stubTimesheetExportUsersByRole{}, &stubTimesheetExportSites{}, file.NewCSVWriter())
		filename, err := uc.ExportIncomeFromTimesheet("individual", "0")

		assert.NoError(t, err)
		assert.NotEmpty(t, filename)
	})

	t.Run("propagates an error from the repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		incomeRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		uc := NewExportIncomeFromTimesheetUsecase(incomeRepo, &stubTimesheetExportStudentLoans{}, &stubTimesheetExportUsersByRole{}, &stubTimesheetExportSites{}, file.NewCSVWriter())
		_, err := uc.ExportIncomeFromTimesheet("individual", "0")

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("resolves month index 0 to the current month", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		record := &models.IncomeFromTimesheet{Income: models.MockIncome}
		record.UserID = bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539").Hex()

		startDate, endDate := models.GetStartDateAndEndDate(time.Now())
		incomeRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return([]*models.IncomeFromTimesheet{record}, nil)

		uc := NewExportIncomeFromTimesheetUsecase(incomeRepo, &stubTimesheetExportStudentLoans{}, &stubTimesheetExportUsersByRole{}, &stubTimesheetExportSites{}, file.NewCSVWriter())
		_, err := uc.ExportIncomeFromTimesheet("individual", "0")

		assert.NoError(t, err)
	})

	t.Run("resolves a non-zero month index to the previous month", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		record := &models.IncomeFromTimesheet{Income: models.MockIncome}
		record.UserID = bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539").Hex()

		startDate, endDate := models.GetStartDateAndEndDate(time.Now().AddDate(0, -1, 0))
		incomeRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return([]*models.IncomeFromTimesheet{record}, nil)

		uc := NewExportIncomeFromTimesheetUsecase(incomeRepo, &stubTimesheetExportStudentLoans{}, &stubTimesheetExportUsersByRole{}, &stubTimesheetExportSites{}, file.NewCSVWriter())
		_, err := uc.ExportIncomeFromTimesheet("individual", "1")

		assert.NoError(t, err)
	})
}

func TestExportIncomeFromTimesheetByStartDateAndEndDate(t *testing.T) {
	t.Run("queries the repository with the given period and writes a CSV", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		startDate := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)

		record := &models.IncomeFromTimesheet{Income: models.MockIncome}
		record.UserID = bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539").Hex()

		incomeRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).
			Return([]*models.IncomeFromTimesheet{record}, nil)

		uc := NewExportIncomeFromTimesheetUsecase(incomeRepo, &stubTimesheetExportStudentLoans{}, &stubTimesheetExportUsersByRole{}, &stubTimesheetExportSites{}, file.NewCSVWriter())
		filename, err := uc.ExportIncomeFromTimesheetByStartDateAndEndDate("individual", startDate, endDate)

		assert.NoError(t, err)
		assert.NotEmpty(t, filename)
	})

	t.Run("propagates an error from the repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)

		incomeRepo.EXPECT().GetAllByRoleStartDateAndEndDate("individual", gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		uc := NewExportIncomeFromTimesheetUsecase(incomeRepo, &stubTimesheetExportStudentLoans{}, &stubTimesheetExportUsersByRole{}, &stubTimesheetExportSites{}, file.NewCSVWriter())
		_, err := uc.ExportIncomeFromTimesheetByStartDateAndEndDate("individual", time.Now(), time.Now())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
