package income

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	iincomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecaseExportIncome(t *testing.T) {
	t.Run("export corporate income success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		year, month := utils.GetYearMonthNow()
		mockRepoIncome := iincomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.MockUserById.ID.Hex(), year, month).Return(&iincomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.MockUserById2.ID.Hex(), year, month).Return(&iincomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().AddExport(gomock.Any()).Return(nil)

		mockRepoUser := userMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetUserByRole("corporate").Return(userMock.MockUsers, nil)

		usecase := NewUsecase(mockRepoIncome, mockRepoUser)
		filename, err := usecase.ExportIncome("corporate")

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
	sum, err := calIncomeSum("100000", "Y")
	assert.NoError(t, err)
	assert.Equal(t, "104000.00", sum.Net)

	sum, err = calIncomeSum("123456", "Y")
	assert.NoError(t, err)
	assert.Equal(t, "128394.24", sum.Net)
}

func TestCalPersonIncomeSum(t *testing.T) {
	sum, err := calIncomeSum("100000", "N")
	assert.NoError(t, err)
	assert.Equal(t, "97000.00", sum.Net)

	sum, err = calIncomeSum("123456", "N")
	assert.NoError(t, err)
	assert.Equal(t, "119752.32", sum.Net)
}

func TestUsecaseAddIncome(t *testing.T) {
	t.Run("when add income success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.MockUser
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockRepoIncome := iincomeMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetUserByID(user.ID.Hex()).Return(&user, nil)
		mockRepoIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(iincomeMock.MockIncome.UserID, year, month).Return(&iincomeMock.MockIncome, errors.New(""))

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.AddIncome(&iincomeMock.MockIncomeReq, &userMock.MockUser)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, iincomeMock.MockIncome.UserID, res.UserID)
	})
}

func TestUsecaseUpdateIncome(t *testing.T) {
	t.Run("when update income success it should return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		user := userMock.MockUser

		mockRepoUser := userMock.NewMockRepository(ctrl)
		incomeMockRepo := iincomeMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetUserByID(user.ID.Hex()).Return(&user, nil)
		incomeMockRepo.EXPECT().UpdateIncome(&iincomeMock.MockIncome).Return(nil)
		incomeMockRepo.EXPECT().GetIncomeByID(iincomeMock.MockIncome.ID.Hex(), userMock.MockUser.ID.Hex()).Return(&iincomeMock.MockIncome, nil)

		uc := NewUsecase(incomeMockRepo, mockRepoUser)
		res, err := uc.UpdateIncome(iincomeMock.MockIncome.ID.Hex(), &iincomeMock.MockIncomeReq, &userMock.MockUser)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, iincomeMock.MockIncome.UserID, res.UserID)
	})
}

func TestUsecaseGetListIncome(t *testing.T) {
	t.Run("when get list income success it should be return list income where status is Y", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := iincomeMock.NewMockRepository(ctrl)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.MockUserById.ID.Hex(), year, month).Return(&iincomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.MockUserById2.ID.Hex(), year, month).Return(&iincomeMock.MockIncome, nil)

		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetUserByRole("corporate").Return(userMock.MockUsers, nil)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeStatusList("corporate")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, iincomeMock.MockIncomeStatusList[0].Status, res[0].Status)

	})
}
func TestUsecaseGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
	t.Run("when get income by user id current month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := iincomeMock.NewMockRepository(ctrl)
		year, month := utils.GetYearMonthNow()
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(iincomeMock.MockIncome.UserID, year, month).Return(&iincomeMock.MockIncome, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAndCurrentMonth(iincomeMock.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, iincomeMock.MockIncome.SubmitDate, res.SubmitDate)
	})
}

func TestSetValueCSV(t *testing.T) {
	assert.Equal(t, `="1"`, setValueCSV("1"))
	assert.Equal(t, `="01"`, setValueCSV("01"))
}
