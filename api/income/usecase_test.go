package income

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecaseExportIncome(t *testing.T) {
	t.Run("export corporate income success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		year, month := utils.GetYearMonthNow()
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().AddExport(gomock.Any()).Return(nil)

		mockRepoUser := userMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		usecase := NewUsecase(mockRepoIncome, mockRepoUser)
		filename, err := usecase.ExportIncome("corporate", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
	t.Run("export corporate income beforeMonth success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		year, month := utils.GetYearMonthNow()
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month-1).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month-1).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().AddExport(gomock.Any()).Return(nil)

		mockRepoUser := userMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		usecase := NewUsecase(mockRepoIncome, mockRepoUser)
		filename, err := usecase.ExportIncome("corporate", "1")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}

func TestCalVAT(t *testing.T) {
	vat, vatf, err := calVAT("100000")
	assert.NoError(t, err)
	assert.Equal(t, "7000.00", vat)
	assert.Equal(t, 7000.00, vatf)

	vat, vatf, err = calVAT("123456")
	assert.NoError(t, err)
	assert.Equal(t, "8641.92", vat)
	assert.Equal(t, 8641.92, vatf)

	vat, vatf, err = calVAT("99999")
	assert.NoError(t, err)
	assert.Equal(t, "6999.93", vat)
	assert.Equal(t, 6999.93, vatf)
}
func TestCalWHT(t *testing.T) {
	wht, whtf, err := calWHT("100000")
	assert.NoError(t, err)
	assert.Equal(t, "3000.00", wht)
	assert.Equal(t, 3000.0, whtf)

	wht, whtf, err = calWHT("123456")
	assert.NoError(t, err)
	assert.Equal(t, "3703.68", wht)
	assert.Equal(t, 3703.68, whtf)

	wht, whtf, err = calWHT("99999")
	assert.NoError(t, err)
	assert.Equal(t, "2999.97", wht)
	assert.Equal(t, 2999.97, whtf)
}
func TestCalCorporateIncomeSum(t *testing.T) {
	sum, err := calIncomeSum("20", "Y", "5000")
	assert.NoError(t, err)
	assert.Equal(t, "104000.00", sum.Net)

	sum, err = calIncomeSum("20", "Y", "2000")
	assert.NoError(t, err)
	assert.Equal(t, "41600.00", sum.Net)
}

func TestCalPersonIncomeSum(t *testing.T) {
	sum, err := calIncomeSum("20", "N", "5000")
	assert.NoError(t, err)
	assert.Equal(t, "97000.00", sum.Net)

	sum, err = calIncomeSum("20", "N", "2000")
	assert.NoError(t, err)
	assert.Equal(t, "38800.00", sum.Net)
}

func TestUsecaseAddIncome(t *testing.T) {
	t.Run("when add income success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(incomeMock.MockIncome.UserID, year, month).Return(&incomeMock.MockIncome, errors.New(""))

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.AddIncome(&incomeMock.MockIncomeReq, &userMock.User)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.UserID, res.UserID)
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
		incomeMockRepo.EXPECT().UpdateIncome(&incomeMock.MockIncome).Return(nil)
		incomeMockRepo.EXPECT().GetIncomeByID(incomeMock.MockIncome.ID.Hex(), userMock.User.ID.Hex()).Return(&incomeMock.MockIncome, nil)

		uc := NewUsecase(incomeMockRepo, mockRepoUser)
		res, err := uc.UpdateIncome(incomeMock.MockIncome.ID.Hex(), &incomeMock.MockIncomeReq, &userMock.User)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.UserID, res.UserID)
	})
}

func TestUsecaseGetListIncome(t *testing.T) {
	t.Run("when get list income success it should be return list income where status is Y", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)

		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeStatusList("corporate")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncomeStatusList[0].Status, res[0].Status)

	})
}
func TestUsecaseGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
	t.Run("when get income by user id current month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(incomeMock.MockIncome.UserID, year, month).Return(&incomeMock.MockIncome, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAndCurrentMonth(incomeMock.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.SubmitDate, res.SubmitDate)
	})
}

func TestSetValueCSV(t *testing.T) {
	assert.Equal(t, `="1"`, setValueCSV("1"))
	assert.Equal(t, `="01"`, setValueCSV("01"))
}
