package income

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/api/income/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
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
		mockRepoIncome.EXPECT().GetStudentLoans()

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
		mockRepoIncome.EXPECT().GetStudentLoans()

		mockRepoUser := userMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		usecase := NewUsecase(mockRepoIncome, mockRepoUser)
		filename, err := usecase.ExportIncome("corporate", "1")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
	t.Run("export individual income current month success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		year, month := utils.GetYearMonthNow()
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeUserByYearMonth(userMock.User2.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().AddExport(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().GetStudentLoans()

		mockRepoUser := userMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetByRole("individual").Return(userMock.Users, nil)

		usecase := NewUsecase(mockRepoIncome, mockRepoUser)
		filename, err := usecase.ExportIncome("individual", "0")

		assert.NoError(t, err)
		assert.NotNil(t, filename)

		// remove file after test
		os.Remove(filename)
	})
}

func TestCSVHeaders(t *testing.T) {
	actual := createHeaders()
	expected := [...]string{"ชื่อ", "ชื่อบัญชี", "เลขบัญชี",
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
	expectedSummaryIncome := "106.70"
	assert.Equal(t, "first last", actual[0])
	assert.Equal(t, "ชื่อ นามสกุล", actual[1])
	assert.Equal(t, expectedAccountNo, actual[2])
	assert.Equal(t, "email@example.com", actual[3])
	assert.Equal(t, expectedNetDailyIncome, actual[4])
	assert.Equal(t, expectedNetSpecialIncome, actual[5])
	assert.Equal(t, expectedWHT, actual[7])
	assert.Equal(t, expectedSummaryIncome, actual[8])
	assert.Equal(t, "note", actual[9])
	assert.Equal(t, "01/12/2022 20:30:00", actual[10])
}

func TestStudentLoanInCSVContent(t *testing.T) {
	loan := models.StudentLoan{
		Fullname: userMock.Admin.BankAccountName,
		Amount:   10,
	}
	actual := createRow(incomeMock.MockIndividualIncome, userMock.IndividualUser1, loan)
	expectedSummaryIncome := "96.70"
	assert.Equal(t, "10.00", actual[6])
	assert.Equal(t, expectedSummaryIncome, actual[8])
}

func TestForeignStudentDoesNotRequireSocialSecuritySoWeUseNegativeStudentLoanToAdjust(t *testing.T) {
	loan := models.StudentLoan{
		Fullname: userMock.Admin.BankAccountName,
		Amount:   -270,
	}
	actual := createRow(incomeMock.MockIndividualIncome, userMock.IndividualUser1, loan)
	expectedSummaryIncome := "376.70"
	assert.Equal(t, "-270.00", actual[6])
	assert.Equal(t, expectedSummaryIncome, actual[8])
}

func TestUseCaseExportIncomeNotExport(t *testing.T) {
	t.Run("export corporate income not export success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoUser := userMock.NewMockRepository(ctrl)
		mockRepoUser.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		year, month := utils.GetYearMonthNow()
		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockRepoIncome.EXPECT().GetIncomeByUserID(userMock.User.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().GetIncomeByUserID(userMock.User2.ID.Hex(), year, month).Return(&incomeMock.MockIncome, nil)
		mockRepoIncome.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().UpdateExportStatus(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().AddExport(gomock.Any()).Return(nil)
		mockRepoIncome.EXPECT().GetStudentLoans()

		usecase := NewUsecase(mockRepoIncome, mockRepoUser)
		filename, err := usecase.ExportIncomeNotExport("corporate")

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
func TestCalWHTCorporate(t *testing.T) {
	wht, whtf, err := calWHTCorporate("100000")
	assert.NoError(t, err)
	assert.Equal(t, "3000.00", wht)
	assert.Equal(t, 3000.00, whtf)

	wht, whtf, err = calWHTCorporate("123456")
	assert.NoError(t, err)
	assert.Equal(t, "3703.68", wht)
	assert.Equal(t, 3703.68, whtf)

	wht, whtf, err = calWHTCorporate("99999")
	assert.NoError(t, err)
	assert.Equal(t, "2999.97", wht)
	assert.Equal(t, 2999.97, whtf)
}
func TestCalCorporateIncomeSum(t *testing.T) {
	sum, err := calIncomeSum("20", "Y", "5000", "corporate")
	assert.NoError(t, err)
	assert.Equal(t, "104000.00", sum.Net)

	sum, err = calIncomeSum("20", "Y", "2000", "corporate")
	assert.NoError(t, err)
	assert.Equal(t, "41600.00", sum.Net)
}

func TestCalPersonIncomeSum(t *testing.T) {
	sum, err := calIncomeSum("20", "N", "5000", "individual")
	assert.NoError(t, err)
	assert.Equal(t, "97000.00", sum.Net)

	sum, err = calIncomeSum("20", "N", "2000", "individual")
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
		res, err := uc.AddIncome(&incomeMock.MockIncomeReq, userMock.User.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, incomeMock.MockIncome.UserID, res.UserID)
		assert.Equal(t, "116400.00", res.NetIncome)
		assert.Equal(t, "97000.00", res.NetDailyIncome)
		assert.Equal(t, "19400.00", res.NetSpecialIncome)
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
