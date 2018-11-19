package income

// import (
// 	"testing"

// 	"github.com/pkg/errors"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/stretchr/testify/mock"
// 	"gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
// 	incomeMocks "gitlab.odds.team/worklog/api.odds-worklog/api/income/mocks"
// 	userMocks "gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
// 	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
// )

// func TestCalVAT(t *testing.T) {
// 	vat, vatf, err := calVAT("100000")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "7000.00", vat)
// 	assert.Equal(t, 7000.00, vatf)

// 	vat, vatf, err = calVAT("123456")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "8641.92", vat)
// 	assert.Equal(t, 8641.92, vatf)

// 	vat, vatf, err = calVAT("99999")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "6999.93", vat)
// 	assert.Equal(t, 6999.93, vatf)
// }

// func TestCalWHT(t *testing.T) {
// 	wht, whtf, err := calWHT("100000")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "3000.00", wht)
// 	assert.Equal(t, 3000.0, whtf)

// 	wht, whtf, err = calWHT("123456")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "3703.68", wht)
// 	assert.Equal(t, 3703.68, whtf)

// 	wht, whtf, err = calWHT("99999")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "2999.97", wht)
// 	assert.Equal(t, 2999.97, whtf)
// }

// func TestCalCorporateIncomeSum(t *testing.T) {
// 	sum, err := calIncomeSum("100000", "Y")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "104000.00", sum.Net)

// 	sum, err = calIncomeSum("123456", "Y")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "128394.24", sum.Net)
// }

// func TestCalPersonIncomeSum(t *testing.T) {
// 	sum, err := calIncomeSum("100000", "N")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "93000.00", sum.Net)

// 	sum, err = calIncomeSum("123456", "N")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "114814.08", sum.Net)
// }

// func TestUsecaseAddIncome(t *testing.T) {
// 	mockUserRepo := new(userMocks.Repository)
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("AddIncome", mock.AnythingOfType("*models.Income")).Return(nil)
// 	year, month := utils.GetYearMonthNow()
// 	mockRepo.On("GetIncomeUserByYearMonth", mocks.MockIncome.UserID, year, month).Return(&mocks.MockIncome, errors.New(""))

// 	uc := NewUsecase(mockRepo, mockUserRepo)
// 	res, err := uc.AddIncome(&mocks.MockIncomeReq, &userMocks.MockUser)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	assert.Equal(t, mocks.MockIncome.UserID, res.UserID)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUsecaseUpdateIncome(t *testing.T) {
// 	mockUserRepo := new(userMocks.Repository)
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("UpdateIncome", &mocks.MockIncome).Return(nil)
// 	mockRepo.On("GetIncomeByID", mocks.MockIncome.ID.Hex(), userMocks.MockUser.ID.Hex()).Return(&mocks.MockIncome, nil)

// 	uc := NewUsecase(mockRepo, mockUserRepo)
// 	res, err := uc.UpdateIncome(mocks.MockIncome.ID.Hex(), &mocks.MockIncomeReq, &userMocks.MockUser)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	assert.Equal(t, mocks.MockIncome.UserID, res.UserID)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUsecaseGetListIncome(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	year, month := utils.GetYearMonthNow()
// 	mockRepo.On("GetIncomeUserByYearMonth", userMocks.MockUserById.ID.Hex(), year, month).Return(&mocks.MockIncome, nil)
// 	mockRepo.On("GetIncomeUserByYearMonth", userMocks.MockUserById2.ID.Hex(), year, month).Return(&mocks.MockIncome, nil)

// 	mockUserRepo := new(userMocks.Repository)
// 	mockUserRepo.On("GetUserByType", "Y").Return(userMocks.MockUsers, nil)

// 	uc := NewUsecase(mockRepo, mockUserRepo)
// 	res, err := uc.GetIncomeStatusList("Y")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	assert.Equal(t, mocks.MockIncomeStatusList[0].Status, res[0].Status)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUsecaseGetIncomeByUserIdAndCurrentMonth(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	year, month := utils.GetYearMonthNow()
// 	mockRepo.On("GetIncomeUserByYearMonth", incomeMocks.MockIncome.UserID, year, month).Return(&mocks.MockIncome, nil)

// 	mockUserRepo := new(userMocks.Repository)

// 	uc := NewUsecase(mockRepo, mockUserRepo)
// 	res, err := uc.GetIncomeByUserIdAndCurrentMonth(mocks.MockIncome.UserID)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	assert.Equal(t, mocks.MockIncome.SubmitDate, res.SubmitDate)
// }

// func TestUsecaseExportIncome(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	year, month := utils.GetYearMonthNow()
// 	mockRepo.On("GetIncomeUserByYearMonth", userMocks.MockUserById.ID.Hex(), year, month).Return(&mocks.MockIncome, nil)
// 	mockRepo.On("GetIncomeUserByYearMonth", userMocks.MockUserById2.ID.Hex(), year, month).Return(&mocks.MockIncome, nil)
// 	mockRepo.On("AddExport", mock.AnythingOfType("*models.Export")).Return(nil)

// 	mockUserRepo := new(userMocks.Repository)
// 	mockUserRepo.On("GetUserByType", "Y").Return(userMocks.MockUsers, nil)

// 	uc := NewUsecase(mockRepo, mockUserRepo)
// 	res, err := uc.ExportIncome("Y")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	mockRepo.AssertExpectations(t)
// 	mockUserRepo.AssertExpectations(t)
// }

// func TestUsecaseDropIncome(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockUserRepo := new(userMocks.Repository)
// 	mockRepo.On("DropIncome").Return(nil)

// 	uc := NewUsecase(mockRepo, mockUserRepo)
// 	err := uc.DropIncome()

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// 	mockUserRepo.AssertExpectations(t)
// }
