package usecases

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestNormalizeTimesheetEventLogLimit(t *testing.T) {
	assert.Equal(t, 100, normalizeTimesheetEventLogLimit(0))
	assert.Equal(t, 100, normalizeTimesheetEventLogLimit(-1))
	assert.Equal(t, 50, normalizeTimesheetEventLogLimit(50))
	assert.Equal(t, 200, normalizeTimesheetEventLogLimit(200))
	assert.Equal(t, 200, normalizeTimesheetEventLogLimit(9999))
}

func TestViewTimesheetEventLogsUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockList := mock_usecases.NewMockForListingTimesheetEventLogs(ctrl)
	mockList.EXPECT().List(100).Return([]*models.TimesheetEventLog{
		{EventType: "timesheet.monthly_summary.published"},
	}, nil)

	uc := NewViewTimesheetEventLogsUsecase(mockList)
	logs, err := uc.List(0)
	assert.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, "timesheet.monthly_summary.published", logs[0].EventType)
}

func TestViewTimesheetEventLogsUsecase_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockList := mock_usecases.NewMockForListingTimesheetEventLogs(ctrl)
	mockList.EXPECT().List(20).Return(nil, errors.New("mongo down"))

	uc := NewViewTimesheetEventLogsUsecase(mockList)
	logs, err := uc.List(20)
	assert.Error(t, err)
	assert.Nil(t, logs)
}
