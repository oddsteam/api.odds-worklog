package income

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	mock_income "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecaseExportIncome(t *testing.T) {
	// t.Run("export corporate income beforeMonth success", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	mockRepoIncome := mockIncomeRepository(ctrl)
	// 	mockRepoIncome.expectGetIncomeUserOfPreviousMonth(userMock.User)
	// 	mockRepoIncome.expectGetIncomeUserOfPreviousMonth(userMock.User2)

	// 	mockRepoUser := userMock.NewMockRepository(ctrl)
	// 	mockRepoUser.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

	// 	usecase := NewUsecase(mockRepoIncome.mock, mockRepoUser)
	// 	filename, err := usecase.ExportIncome("corporate", "1")

	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, filename)

	// 	// remove file after test
	// 	os.Remove(filename)
	// })
	t.Run("export individual income current month success (new)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		incomes := []*models.Income{
			&incomeMock.MockIncome,
			&incomeMock.MockIncome2,
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.expectNumberOfExportStatuses(2)
		mockRepoIncome.expectGetAllIncomeOfCurrentMonthByRole(incomes)

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncome("individual", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})

	t.Run("export individual income (new) includes income from friendslog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		incomes := []*models.Income{
			&incomeMock.MockIncome,
			&incomeMock.MockIncome2,
			{ID: bson.ObjectIdHex("5bd1fda30fd2df2a3e41e571"), Role: "individual", WorkDate: "20", DailyRate: 750},
		}
		mockRepoIncome := mockIncomeRepository(ctrl)
		mockRepoIncome.expectNumberOfExportStatuses(3)
		mockRepoIncome.expectGetAllIncomeOfCurrentMonthByRole(incomes)

		usecase := NewUsecase(mockRepoIncome.mock, userMock.NewMockRepository(ctrl))
		filename, err := usecase.ExportIncome("individual", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}

func TestCSVHeaders(t *testing.T) {
	actual := createHeaders()
	expected := [...]string{"Vendor Code", "ชื่อบัญชี", "Payment method", "เลขบัญชี", "ชื่อ", "เลขบัตรประชาชน",
		"อีเมล", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ",
		"หัก ณ ที่จ่าย", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก",
	}
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

func TestCSVContentForIndividual(t *testing.T) {
	actual := createRow(incomeMock.MockIndividualIncome, userMock.IndividualUser1, models.StudentLoan{})
	expectedAccountNo := `="0531231231"`
	expectedNetDailyIncome := "97.00"
	expectedNetSpecialIncome := "9.70"
	expectedWHT := "3.30"
	expectedTransferAmount := "106.70"
	expected := [...]string{"", "ชื่อ นามสกุล", "", expectedAccountNo, "first last", "ThaiCitizenID",
		"email@example.com", expectedNetDailyIncome, expectedNetSpecialIncome,
	}
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
	assert.Equal(t, expectedWHT, actual[WITHHOLDING_TAX_INDEX])
	assert.Equal(t, expectedTransferAmount, actual[TRANSFER_AMOUNT_INDEX])
	assert.Equal(t, "note", actual[NOTE_INDEX])
	assert.Equal(t, "01/12/2022 20:30:00", actual[SUBMIT_DATE_INDEX])
}

func TestStudentLoanInCSVContent(t *testing.T) {
	loan := models.StudentLoan{
		Fullname: userMock.Admin.BankAccountName,
		Amount:   10,
	}
	actual := createRow(incomeMock.MockIndividualIncome, userMock.IndividualUser1, loan)
	expectedTransferAmount := "96.70"
	assert.Equal(t, "10.00", actual[LOAN_DEDUCTION_INDEX])
	assert.Equal(t, expectedTransferAmount, actual[TRANSFER_AMOUNT_INDEX])
}

func TestForeignStudentDoesNotRequireSocialSecuritySoWeUseNegativeStudentLoanToAdjust(t *testing.T) {
	loan := models.StudentLoan{
		Fullname: userMock.Admin.BankAccountName,
		Amount:   -270,
	}
	actual := createRow(incomeMock.MockIndividualIncome, userMock.IndividualUser1, loan)
	expectedTransferAmount := "376.70"
	assert.Equal(t, "-270.00", actual[LOAN_DEDUCTION_INDEX])
	assert.Equal(t, expectedTransferAmount, actual[TRANSFER_AMOUNT_INDEX])
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
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(incomeMock.MockIncome.UserID, year, month).Return(&incomeMock.MockIncome, errors.New(""))

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.AddIncome(&incomeMock.MockIncomeReq, userMock.User.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.UserID, res.UserID)
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
		incomeMockRepo.EXPECT().UpdateIncome(&incomeMock.MockIncome).Return(nil)
		incomeMockRepo.EXPECT().GetIncomeByID(incomeMock.MockIncome.ID.Hex(), userMock.User.ID.Hex()).Return(&incomeMock.MockIncome, nil)

		uc := NewUsecase(incomeMockRepo, mockRepoUser)
		res, err := uc.UpdateIncome(incomeMock.MockIncome.ID.Hex(), &incomeMock.MockIncomeReq, userMock.User.ID.Hex())

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
		res, err := uc.GetIncomeStatusList("corporate", false)
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
func TestUsecaseGetIncomeByUserIdAndAllMonth(t *testing.T) {
	t.Run("when get income by user id all month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserIdAllMonth(incomeMock.MockIncome.UserID).Return(incomeMock.MockIncomeList, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAllMonth(incomeMock.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.SubmitDate, res[0].SubmitDate)
		assert.Equal(t, "50440.00", res[1].NetIncome)
	})
}

func TestUsecaseGetIncomeByUserIdAndAllMonthCaseNoNetSpecialIncome(t *testing.T) {
	t.Run("when get income by user id all month success it should be return income model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserIdAllMonth(incomeMock.MockIncome.UserID).Return(incomeMock.MockIncomeListNoNetSpecialIncome, nil)
		mockUserRepo := userMock.NewMockRepository(ctrl)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetIncomeByUserIdAllMonth(incomeMock.MockIncome.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "00.00", res[0].NetIncome)
		assert.Equal(t, "50440.00", res[1].NetIncome)
	})
}

type MockIncomeRepository struct {
	mock *mock_income.MockRepository
}

func (m *MockIncomeRepository) expectUpdateExportStatus() {
	m.mock.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
}

func (m *MockIncomeRepository) expectGetIncomeUserOfCurrentMonth(u models.User) {
	year, month := utils.GetYearMonthNow()
	m.mock.EXPECT().GetIncomeUserByYearMonth(u.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
}

func (m *MockIncomeRepository) expectGetIncomeUserOfPreviousMonth(u models.User) {
	year, month := utils.GetYearMonthNow()
	m.mock.EXPECT().GetIncomeUserByYearMonth(u.ID.Hex(), year, month-1).Return(&incomeMock.MockIncome, nil)
}

func (m *MockIncomeRepository) expectGetAllIncomeOfCurrentMonth(incomes []*models.Income) {
	startDate, endDate := utils.GetStartDateAndEndDate(time.Now())
	m.mock.EXPECT().GetAllIncomeByStartDateAndEndDate(
		gomock.Any(), startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) expectGetAllIncomeOfCurrentMonthByRole(incomes []*models.Income) {
	startDate, endDate := utils.GetStartDateAndEndDate(time.Now())
	m.mock.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		gomock.Any(), startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) expectNumberOfExportStatuses(times int) {
	for i := 0; i < times; i++ {
		m.expectUpdateExportStatus()
	}
}

func mockIncomeRepository(ctrl *gomock.Controller) *MockIncomeRepository {
	mockRepoIncome := MockIncomeRepository{incomeMock.NewMockRepository(ctrl)}
	mockRepoIncome.mock.EXPECT().AddExport(gomock.Any()).Return(nil)
	mockRepoIncome.mock.EXPECT().GetStudentLoans()
	return &mockRepoIncome
}
