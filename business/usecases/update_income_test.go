package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

// incomeSubmittedIn returns a copy of the mock income stamped with a fixed submit date, so a
// test can assert which period the income_from_timesheet lookup uses.
func incomeSubmittedIn(year int, month time.Month) *models.Income {
	income := models.MockIncome
	income.SubmitDate = time.Date(year, month, 15, 9, 0, 0, 0, time.UTC)
	return &income
}

// existingTimesheetRecord is the record the timesheet consumer would already have written for
// the period, complete with its own id and per-site breakdown.
func existingTimesheetRecord() *models.IncomeFromTimesheet {
	income := models.MockIncome
	income.ID = bsonutil.MustObjectIDFromHex("6bd1fda30fd2df2a3e41e569")
	income.WorkDate = "12.5"
	income.WorkingHours = "2"
	return &models.IncomeFromTimesheet{
		Income: income,
		Sites: []models.SiteWork{
			{ClientSite: "SITE-A", CustomerName: "Site A Customer", WorkingDays: 10, OvertimeDays: 1},
			{ClientSite: "SITE-B", CustomerName: "Site B Customer", WorkingDays: 2.5, OvertimeDays: 1},
		},
	}
}

func TestUsecaseUpdateIncome(t *testing.T) {
	t.Run("when update income success it should return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User

		mockUserRepo := mock_usecases.NewMockForGettingUserByID(ctrl)
		mockIncomeRepo := mock_usecases.NewMockForUpdatingUserIncome(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockIncomeRepo.EXPECT().GetIncomeByID(models.MockIncome.ID.Hex(), user.ID.Hex()).Return(incomeSubmittedIn(2026, time.June), nil)
		mockIncomeRepo.EXPECT().UpdateIncome(gomock.Any()).Return(nil)
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		mockTimesheetRepo.EXPECT().Add(gomock.Any()).Return(nil)

		uc := NewUpdateIncomeUsecase(mockIncomeRepo, mockUserRepo, mockTimesheetRepo)
		res, err := uc.UpdateIncome(models.MockIncome.ID.Hex(), &models.MockIncomeReq, user.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncome.UserID, res.UserID)
	})

	t.Run("looks up the income_from_timesheet record by the period the edited income belongs to", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User

		mockUserRepo := mock_usecases.NewMockForGettingUserByID(ctrl)
		mockIncomeRepo := mock_usecases.NewMockForUpdatingUserIncome(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockIncomeRepo.EXPECT().GetIncomeByID(gomock.Any(), gomock.Any()).Return(incomeSubmittedIn(2026, time.June), nil)
		mockIncomeRepo.EXPECT().UpdateIncome(gomock.Any()).Return(nil)
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(models.MockIncome.UserID, 2026, time.June).Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		mockTimesheetRepo.EXPECT().Add(gomock.Any()).Return(nil)

		uc := NewUpdateIncomeUsecase(mockIncomeRepo, mockUserRepo, mockTimesheetRepo)
		_, err := uc.UpdateIncome(models.MockIncome.ID.Hex(), &models.MockIncomeReq, user.ID.Hex())

		assert.NoError(t, err)
	})

	t.Run("overwrites an existing income_from_timesheet record but keeps its sites", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User

		mockUserRepo := mock_usecases.NewMockForGettingUserByID(ctrl)
		mockIncomeRepo := mock_usecases.NewMockForUpdatingUserIncome(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockIncomeRepo.EXPECT().GetIncomeByID(gomock.Any(), gomock.Any()).Return(incomeSubmittedIn(2026, time.June), nil)
		mockIncomeRepo.EXPECT().UpdateIncome(gomock.Any()).Return(nil)

		existing := existingTimesheetRecord()
		existingID := existing.ID
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(models.MockIncome.UserID, 2026, time.June).Return(existing, nil)

		var saved *models.IncomeFromTimesheet
		mockTimesheetRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(rec *models.IncomeFromTimesheet) error {
			saved = rec
			return nil
		})

		uc := NewUpdateIncomeUsecase(mockIncomeRepo, mockUserRepo, mockTimesheetRepo)
		res, err := uc.UpdateIncome(models.MockIncome.ID.Hex(), &models.MockIncomeReq, user.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, saved)
		assert.Equal(t, existingID, saved.ID, "must keep the income_from_timesheet record's own id")
		assert.Equal(t, res.WorkDate, saved.WorkDate)
		assert.Equal(t, res.WorkingHours, saved.WorkingHours)
		assert.Equal(t, res.NetIncome, saved.NetIncome)
		assert.Equal(t, existingTimesheetRecord().Sites, saved.Sites)
	})

	t.Run("returns an error when writing the income_from_timesheet record fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User

		mockUserRepo := mock_usecases.NewMockForGettingUserByID(ctrl)
		mockIncomeRepo := mock_usecases.NewMockForUpdatingUserIncome(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockIncomeRepo.EXPECT().GetIncomeByID(gomock.Any(), gomock.Any()).Return(incomeSubmittedIn(2026, time.June), nil)
		mockIncomeRepo.EXPECT().UpdateIncome(gomock.Any()).Return(nil)
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		uc := NewUpdateIncomeUsecase(mockIncomeRepo, mockUserRepo, mockTimesheetRepo)
		res, err := uc.UpdateIncome(models.MockIncome.ID.Hex(), &models.MockIncomeReq, user.ID.Hex())

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)
	})
}
