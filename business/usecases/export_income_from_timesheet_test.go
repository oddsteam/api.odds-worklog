package usecases

import (
	"testing"

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
}
