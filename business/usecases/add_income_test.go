package usecases

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/business/models/mock"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestUsecaseAddIncome(t *testing.T) {
	t.Run("when a user add income success it should be return income model to show on the screen", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		year, month := models.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(models.MockIncome.UserID, year, month).Return(&models.MockIncome, errors.New(""))
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		mockTimesheetRepo.EXPECT().Add(gomock.Any()).Return(nil)

		uc := NewAddIncomeUsecase(mockRepoIncome, mockUserRepo, mockTimesheetRepo)
		res, err := uc.AddIncome(&models.MockIncomeReq, userMock.User.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, models.MockIncome.UserID, res.UserID)
		assert.Equal(t, "114000.00", res.NetIncome)
		assert.Equal(t, "95000.00", res.NetDailyIncome)
		assert.Equal(t, "19000.00", res.NetSpecialIncome)
		assert.Equal(t, "", res.VAT)
		assert.Equal(t, "6000.00", res.WHT)
		assert.Equal(t, 0.05, res.WHTRate)
	})

	t.Run("saves the note the user typed, into both income and income_from_timesheet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		year, month := models.GetYearMonthNow()
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(user.ID.Hex(), year, month).Return(nil, errors.New("not found"))
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), year, month).Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)

		var saved *models.IncomeFromTimesheet
		mockTimesheetRepo.EXPECT().Add(gomock.Any()).DoAndReturn(func(rec *models.IncomeFromTimesheet) error {
			saved = rec
			return nil
		})

		req := models.MockIncomeReq
		req.Note = "ลาป่วย 2 วัน"

		uc := NewAddIncomeUsecase(mockRepoIncome, mockUserRepo, mockTimesheetRepo)
		res, err := uc.AddIncome(&req, user.ID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, "ลาป่วย 2 วัน", res.Note)
		assert.Equal(t, "ลาป่วย 2 วัน", saved.Note)
	})

	t.Run("creates an income_from_timesheet record for the period when none exists yet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		year, month := models.GetYearMonthNow()
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(user.ID.Hex(), year, month).Return(nil, errors.New("not found"))
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), year, month).Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)

		var saved *models.IncomeFromTimesheet
		mockTimesheetRepo.EXPECT().Add(gomock.Any()).DoAndReturn(func(rec *models.IncomeFromTimesheet) error {
			saved = rec
			return nil
		})

		uc := NewAddIncomeUsecase(mockRepoIncome, mockUserRepo, mockTimesheetRepo)
		res, err := uc.AddIncome(&models.MockIncomeReq, user.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, saved)
		assert.Equal(t, res.UserID, saved.UserID)
		assert.Equal(t, res.WorkDate, saved.WorkDate)
		assert.Equal(t, res.WorkingHours, saved.WorkingHours)
		assert.Equal(t, res.NetIncome, saved.NetIncome)
		assert.Empty(t, saved.Sites)
	})

	t.Run("overwrites an existing income_from_timesheet record but keeps its sites", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		year, month := models.GetYearMonthNow()
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(user.ID.Hex(), year, month).Return(nil, errors.New("not found"))
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)

		existing := existingTimesheetRecord()
		existingID := existing.ID
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), year, month).Return(existing, nil)

		var saved *models.IncomeFromTimesheet
		mockTimesheetRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(rec *models.IncomeFromTimesheet) error {
			saved = rec
			return nil
		})

		uc := NewAddIncomeUsecase(mockRepoIncome, mockUserRepo, mockTimesheetRepo)
		res, err := uc.AddIncome(&models.MockIncomeReq, user.ID.Hex())

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
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockTimesheetRepo := mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl)
		year, month := models.GetYearMonthNow()
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(user.ID.Hex(), year, month).Return(nil, errors.New("not found"))
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		mockTimesheetRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), year, month).Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
		mockTimesheetRepo.EXPECT().Add(gomock.Any()).Return(assert.AnError)

		uc := NewAddIncomeUsecase(mockRepoIncome, mockUserRepo, mockTimesheetRepo)
		res, err := uc.AddIncome(&models.MockIncomeReq, user.ID.Hex())

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)
	})
}
