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
		ID:    bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539"),
		Email: "test@abc.com",
		Role:  "individual",
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

func TestSyncIncomeForTimesheet(t *testing.T) {
	t.Run("creates a new income_for_timesheet record when none exists, summing days across sites", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(nil, ErrIncomeForTimesheetNotFoundForPeriod)
		incomeRepo.EXPECT().Add(gomock.Any()).DoAndReturn(func(record *models.IncomeForTimesheet) error {
			assert.Equal(t, "12.50", record.WorkDate)
			assert.Equal(t, "2.00", record.WorkingHours)
			assert.Equal(t, "0", record.SpecialIncome)
			assert.Equal(t, []models.SiteWork{
				{ClientSite: "SITE-A", CustomerName: "Site A Customer", WorkingDays: 10, OvertimeDays: 1},
				{ClientSite: "SITE-B", CustomerName: "Site B Customer", WorkingDays: 2.5, OvertimeDays: 1},
			}, record.Sites)
			return nil
		})

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(evt)

		assert.NoError(t, err)
	})

	t.Run("updates the existing record for the month and preserves its note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()
		existing := &models.IncomeForTimesheet{Income: models.MockIncome}
		existing.Note = "existing remark"

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(existing, nil)
		incomeRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(record *models.IncomeForTimesheet) error {
			assert.Equal(t, "existing remark", record.Note)
			assert.Equal(t, "12.50", record.WorkDate)
			assert.Equal(t, "2.00", record.WorkingHours)
			return nil
		})

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(evt)

		assert.NoError(t, err)
	})

	t.Run("returns ErrTimesheetUserNotFound when GetByEmail returns it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(nil, ErrTimesheetUserNotFound)

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(timesheetSyncEvent())

		assert.ErrorIs(t, err, ErrTimesheetUserNotFound)
	})

	t.Run("propagates a real error from GetByEmail instead of masking it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(nil, assert.AnError)
		incomeRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		incomeRepo.EXPECT().Add(gomock.Any()).Times(0)
		incomeRepo.EXPECT().Update(gomock.Any()).Times(0)

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(timesheetSyncEvent())

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("propagates a real error from GetByUserYearMonth instead of treating it as not-found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		user := timesheetSyncUser()

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(nil, assert.AnError)
		incomeRepo.EXPECT().Add(gomock.Any()).Times(0)
		incomeRepo.EXPECT().Update(gomock.Any()).Times(0)

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(timesheetSyncEvent())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
