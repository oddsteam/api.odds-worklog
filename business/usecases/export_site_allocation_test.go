package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

type stubSiteAllocationWriter struct {
	name        string
	allocations []models.SiteAllocation
	err         error
}

func (s *stubSiteAllocationWriter) WriteFile(name string, allocations []models.SiteAllocation) (string, error) {
	s.name = name
	s.allocations = allocations
	if s.err != nil {
		return "", s.err
	}
	return "files/site_allocation_stub.csv", nil
}

func givenTimesheetRecords() []*models.IncomeFromTimesheet {
	return []*models.IncomeFromTimesheet{
		{
			Income: models.Income{DailyIncomeBeforeTax: "1000.00"},
			Sites: []models.SiteWork{
				{ClientSite: "SCB", WorkingDays: 15},
				{ClientSite: "KBANK", WorkingDays: 5},
			},
		},
	}
}

func TestExportSiteAllocationWritesTheBreakdownOfTheRequestedMonth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)
	exportLog := mock_usecases.NewMockForControllingIncomeData(ctrl)
	writer := &stubSiteAllocationWriter{}

	startDate, endDate := models.GetStartDateAndEndDate(time.Now())
	reader.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).Return(givenTimesheetRecords(), nil)
	exportLog.EXPECT().AddExport(gomock.Any()).Return(nil)

	u := NewExportSiteAllocationUsecase(reader, exportLog, writer)
	filename, err := u.ExportSiteAllocation("individual", "0")

	assert.NoError(t, err)
	assert.Equal(t, "files/site_allocation_stub.csv", filename)
	assert.Equal(t, []models.SiteAllocation{
		{Site: "SCB", Amount: 750, Percent: 75},
		{Site: "KBANK", Amount: 250, Percent: 25},
	}, writer.allocations)
}

func TestExportSiteAllocationForAPeriodSpansTheWholeRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)
	exportLog := mock_usecases.NewMockForControllingIncomeData(ctrl)
	writer := &stubSiteAllocationWriter{}

	startDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)
	reader.EXPECT().GetAllByRoleStartDateAndEndDate("individual", startDate, endDate).Return(givenTimesheetRecords(), nil)
	exportLog.EXPECT().AddExport(gomock.Any()).Return(nil)

	u := NewExportSiteAllocationUsecase(reader, exportLog, writer)
	_, err := u.ExportSiteAllocationByStartDateAndEndDate("individual", startDate, endDate)

	assert.NoError(t, err)
	assert.Len(t, writer.allocations, 2)
}

func TestExportSiteAllocationDoesNotLogAnExportThatFailedToWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)
	exportLog := mock_usecases.NewMockForControllingIncomeData(ctrl)
	writer := &stubSiteAllocationWriter{err: errors.New("no data for export to CSV file")}

	reader.EXPECT().GetAllByRoleStartDateAndEndDate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	u := NewExportSiteAllocationUsecase(reader, exportLog, writer)
	_, err := u.ExportSiteAllocation("individual", "0")

	assert.Error(t, err)
}

func TestExportSiteAllocationSurfacesReadFailures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mock_usecases.NewMockForGettingIncomeFromTimesheetInTheMonth(ctrl)
	exportLog := mock_usecases.NewMockForControllingIncomeData(ctrl)
	writer := &stubSiteAllocationWriter{}

	readErr := errors.New("mongo is down")
	reader.EXPECT().GetAllByRoleStartDateAndEndDate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, readErr)

	u := NewExportSiteAllocationUsecase(reader, exportLog, writer)
	_, err := u.ExportSiteAllocation("individual", "0")

	assert.ErrorIs(t, err, readErr)
	assert.Nil(t, writer.allocations)
}
