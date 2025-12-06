package income

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/entity"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/entity/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecaseExportIncome(t *testing.T) {
	t.Run("export individual income current month success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		incomes := []*models.Income{
			&entity.MockIncome,
			&entity.MockIncome2,
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.expectGetAllIncomeOfCurrentMonthByRole(incomes, time.Now())

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
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
		mockRepoIncome.expectGetAllIncomeOfCurrentMonthByRole(incomes, time.Now())

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
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
		mockRepoIncome.expectGetAllIncomeOfPreviousMonthByRole(incomes)

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
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

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
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

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
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

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncomeSAPByStartDateAndEndDate("corporate", startDate, endDate, dateEff)

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}

func TestUsecaseAddIncome(t *testing.T) {
	t.Run("when a user add income success it should be return income model to show on the screen", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(entity.MockIncome.UserID, year, month).Return(&entity.MockIncome, errors.New(""))

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.AddIncome(&entity.MockIncomeReq, userMock.User.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.MockIncome.UserID, res.UserID)
		assert.Equal(t, "116400.00", res.NetIncome)
		assert.Equal(t, "97000.00", res.NetDailyIncome)
		assert.Equal(t, "19400.00", res.NetSpecialIncome)
		assert.Equal(t, "", res.VAT)
		assert.Equal(t, "3600.00", res.WHT)
	})
}

func TestUsecaseUpdateIncome(t *testing.T) {
	t.Run("when update income success it should return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User

		mockRepoUser := userMock.NewMockRepository(ctrl)
		incomeMockRepo := incomeMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		incomeMockRepo.EXPECT().UpdateIncome(&entity.MockIncome).Return(nil)
		incomeMockRepo.EXPECT().GetIncomeByID(entity.MockIncome.ID.Hex(), userMock.User.ID.Hex()).Return(&entity.MockIncome, nil)

		uc := NewUsecase(incomeMockRepo, mockRepoUser)
		res, err := uc.UpdateIncome(entity.MockIncome.ID.Hex(), &entity.MockIncomeReq, userMock.User.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.MockIncome.UserID, res.UserID)
	})
}

func TestUsecaseGetListIncome(t *testing.T) {
	t.Run("when get list income success it should be return list income where status is Y", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month).Return(&entity.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month).Return(&entity.MockIncome, nil)

		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeStatusList("corporate", false)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.MockIncomeStatusList[0].Status, res[0].Status)

	})
}
func TestUsecaseGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
	t.Run("when get income by user id current month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(entity.MockIncome.UserID, year, month).Return(&entity.MockIncome, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAndCurrentMonth(entity.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.MockIncome.SubmitDate, res.SubmitDate)
	})
}
func TestUsecaseGetIncomeByUserIdAndAllMonth(t *testing.T) {
	t.Run("when get income by user id all month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserIdAllMonth(entity.MockIncome.UserID).Return(entity.MockIncomeList, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAllMonth(entity.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.MockIncome.SubmitDate, res[0].SubmitDate)
		assert.Equal(t, "50440.00", res[1].NetIncome)
	})
}

func TestUsecaseGetIncomeByUserIdAndAllMonthCaseNoNetSpecialIncome(t *testing.T) {
	t.Run("when get income by user id all month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserIdAllMonth(entity.MockIncome.UserID).Return(entity.MockIncomeListNoNetSpecialIncome, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAllMonth(entity.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "00.00", res[0].NetIncome)
		assert.Equal(t, "50440.00", res[1].NetIncome)
	})
}

type MockIncomeRepository struct {
	mock *incomeMock.MockRepository
}

func (m *MockIncomeRepository) expectGetAllIncomeOfPreviousMonthByRole(incomes []*models.Income) {
	previousMonth := time.Now().AddDate(0, -1, 0)
	m.expectGetAllIncomeOfCurrentMonthByRole(incomes, previousMonth)
}

func (m *MockIncomeRepository) expectGetAllIncomeOfCurrentMonthByRole(incomes []*models.Income, now time.Time) {
	startDate, endDate := utils.GetStartDateAndEndDate(now)
	m.mock.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		gomock.Any(), startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) GetAllIncomeByRoleStartDateAndEndDate(incomes []*models.Income, role string, startDate, endDate time.Time) {
	m.mock.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		role, startDate, endDate).Return(incomes, nil)
}

func mockIncomeRepository(ctrl *gomock.Controller) *MockIncomeRepository {
	mockRepoIncome := MockIncomeRepository{incomeMock.NewMockRepository(ctrl)}
	mockRepoIncome.mock.EXPECT().AddExport(gomock.Any()).Return(nil)
	mockRepoIncome.mock.EXPECT().GetStudentLoans()
	return &mockRepoIncome
}

func deepClone(income *models.Income) *models.Income {
	b, _ := json.Marshal(income)
	var i models.Income
	json.Unmarshal(b, &i)
	return &i

}
