package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

func timesheetSyncUser() models.User {
	return models.User{
		ID:          bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539"),
		Email:       "test@abc.com",
		Role:        "individual",
		DailyIncome: "800",
	}
}

func timesheetSyncEvent() models.TimesheetMonthlySummaryEvent {
	return models.TimesheetMonthlySummaryEvent{
		EventType: "timesheet.monthly_summary",
		Year:      2026,
		Month:     6,
		SummaryAt: time.Date(2026, 7, 10, 15, 31, 10, 0, time.UTC),
		Employee:  models.TimesheetEmployee{Email: "test@abc.com", EnglishName: "Tester Super"},
		Sites: []models.TimesheetSiteSummary{
			{ClientSite: "SITE-A", CustomerName: "Site A Customer", WorkingDays: 10, OvertimeDays: 1},
			{ClientSite: "SITE-B", CustomerName: "Site B Customer", WorkingDays: 2.5, OvertimeDays: 1},
		},
	}
}

func TestSyncIncomeFromTimesheet(t *testing.T) {
	t.Run("saves the raw event log before doing anything else", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		eventLogRepo := mock_usecases.NewMockForLoggingTimesheetEvent(ctrl)

		evt := timesheetSyncEvent()

		eventLogRepo.EXPECT().Save(evt).Return(assert.AnError)
		userRepo.EXPECT().GetByEmail(gomock.Any()).Times(0)
		incomeRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		uc := NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo, eventLogRepo, nil)
		err := uc.SyncFromEvent(evt)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("creates a new income_from_timesheet record when none exists, summing days across sites", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		eventLogRepo := mock_usecases.NewMockForLoggingTimesheetEvent(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()

		eventLogRepo.EXPECT().Save(evt).Return(nil)
		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		incomeRepo.EXPECT().Add(gomock.Any()).DoAndReturn(func(record *models.IncomeFromTimesheet) error {
			assert.Equal(t, "12.50", record.WorkDate)
			assert.Equal(t, "16.00", record.WorkingHours)
			assert.Equal(t, "100.00", record.SpecialIncome)
			assert.Equal(t, []models.SiteWork{
				{ClientSite: "SITE-A", CustomerName: "Site A Customer", WorkingDays: 10, OvertimeDays: 1},
				{ClientSite: "SITE-B", CustomerName: "Site B Customer", WorkingDays: 2.5, OvertimeDays: 1},
			}, record.Sites)
			return nil
		})

		uc := NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo, eventLogRepo, nil)
		err := uc.SyncFromEvent(evt)

		assert.NoError(t, err)
	})

	t.Run("updates the existing record for the month and preserves its note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		eventLogRepo := mock_usecases.NewMockForLoggingTimesheetEvent(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()
		existing := &models.IncomeFromTimesheet{Income: models.MockIncome}
		existing.Note = "existing remark"

		eventLogRepo.EXPECT().Save(evt).Return(nil)
		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(existing, nil)
		incomeRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(record *models.IncomeFromTimesheet) error {
			assert.Equal(t, "existing remark", record.Note)
			assert.Equal(t, "12.50", record.WorkDate)
			assert.Equal(t, "16.00", record.WorkingHours)
			return nil
		})

		uc := NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo, eventLogRepo, nil)
		err := uc.SyncFromEvent(evt)

		assert.NoError(t, err)
	})

	t.Run("returns ErrTimesheetUserNotFound when GetByEmail returns it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		eventLogRepo := mock_usecases.NewMockForLoggingTimesheetEvent(ctrl)

		evt := timesheetSyncEvent()
		eventLogRepo.EXPECT().Save(evt).Return(nil)
		userRepo.EXPECT().GetByEmail("test@abc.com").Return(nil, ErrTimesheetUserNotFound)

		uc := NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo, eventLogRepo, nil)
		err := uc.SyncFromEvent(evt)

		assert.ErrorIs(t, err, ErrTimesheetUserNotFound)
	})

	t.Run("propagates a real error from GetByEmail instead of masking it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		eventLogRepo := mock_usecases.NewMockForLoggingTimesheetEvent(ctrl)

		evt := timesheetSyncEvent()
		eventLogRepo.EXPECT().Save(evt).Return(nil)
		userRepo.EXPECT().GetByEmail("test@abc.com").Return(nil, assert.AnError)
		incomeRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		incomeRepo.EXPECT().Add(gomock.Any()).Times(0)
		incomeRepo.EXPECT().Update(gomock.Any()).Times(0)

		uc := NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo, eventLogRepo, nil)
		err := uc.SyncFromEvent(evt)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("propagates a real error from GetByUserYearMonth instead of treating it as not-found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		eventLogRepo := mock_usecases.NewMockForLoggingTimesheetEvent(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()
		eventLogRepo.EXPECT().Save(evt).Return(nil)
		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(nil, assert.AnError)
		incomeRepo.EXPECT().Add(gomock.Any()).Times(0)
		incomeRepo.EXPECT().Update(gomock.Any()).Times(0)

		uc := NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo, eventLogRepo, nil)
		err := uc.SyncFromEvent(evt)

		assert.ErrorIs(t, err, assert.AnError)
	})
}
