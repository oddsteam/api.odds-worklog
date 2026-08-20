package usecases

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
)

func CreateExportIncomeUsecaseWithMock(t *testing.T) (ForUsingExportIncome, *gomock.Controller, *MockIncomeRepository) {
	ctrl := gomock.NewController(t)
	mockRepoIncome := mockIncomeRepository(ctrl)

	mockRepoIncome.mockGettingUsersByRole.EXPECT().GetByRole(gomock.Any()).Return([]*models.User{}, nil).AnyTimes()
	mockSites := &stubListingSites{}

	usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, mockRepoIncome.mockSapExportFailure, file.NewCSVWriter(), file.NewSAPWriter(), file.NewPeakCSVWriter(), mockRepoIncome.mockRead, mockRepoIncome.mockGettingUsersByRole, mockSites)
	return usecase, ctrl, mockRepoIncome
}

type stubListingSites struct{}

func (s *stubListingSites) GetSiteGroup() ([]*models.Site, error) {
	return []*models.Site{}, nil
}

func CreateAddIncomeUsecaseWithMock(mockRepoIncome *MockIncomeRepository) ForUsingAddIncome {
	usecase := NewAddIncomeUsecase(mockRepoIncome.mockControllingUserIncome, mockRepoIncome.mockGettingUserByID, mockRepoIncome.mockIncomeFromTimesheet)
	return usecase
}

func CreateGetIncomeUsecaseWithMock(mockRepoIncome *MockIncomeRepository) ForUsingGetIncome {
	usecase := NewGetIncomeUsecase(mockRepoIncome.mockReadingUserIncome)
	return usecase
}

func CreateUpdateIncomeUsecaseWithMock(mockRepoIncome *MockIncomeRepository) ForUsingUpdateIncome {
	usecase := NewUpdateIncomeUsecase(mockRepoIncome.mockUpdatingUserIncome, mockRepoIncome.mockGettingUserByID, mockRepoIncome.mockIncomeFromTimesheet)
	return usecase
}

func CreateListIncomeStatusUsecaseWithMock(mockRepoIncome *MockIncomeRepository) ForUsingListIncomeStatus {
	usecase := NewListIncomeStatusUsecase(mockRepoIncome.mockReadingUserIncome, mockRepoIncome.mockGettingUsersByRole)
	return usecase
}

func mockIncomeRepository(ctrl *gomock.Controller) *MockIncomeRepository {
	mockRepoIncome := MockIncomeRepository{
		mockGettingUserByID:       mock_usecases.NewMockForGettingUserByID(ctrl),
		mockControllingUserIncome: mock_usecases.NewMockForControllingUserIncome(ctrl),
		mockReadingUserIncome:     mock_usecases.NewMockForReadingUserIncome(ctrl),
		mockUpdatingUserIncome:    mock_usecases.NewMockForUpdatingUserIncome(ctrl),
		mockGettingUsersByRole:    mock_usecases.NewMockForGettingUsersByRole(ctrl),
		mockRead:                  mock_usecases.NewMockForGettingIncomeData(ctrl),
		mockWrite:                 mock_usecases.NewMockForControllingIncomeData(ctrl),
		mockSapExportFailure:      mock_usecases.NewMockForLoggingSAPExportFailure(ctrl),
		mockIncomeFromTimesheet:   mock_usecases.NewMockForGettingIncomeFromTimesheet(ctrl),
	}
	return &mockRepoIncome
}

type MockIncomeRepository struct {
	mockGettingUserByID       *mock_usecases.MockForGettingUserByID
	mockControllingUserIncome *mock_usecases.MockForControllingUserIncome
	mockReadingUserIncome     *mock_usecases.MockForReadingUserIncome
	mockUpdatingUserIncome    *mock_usecases.MockForUpdatingUserIncome
	mockGettingUsersByRole    *mock_usecases.MockForGettingUsersByRole
	mockRead                  *mock_usecases.MockForGettingIncomeData
	mockWrite                 *mock_usecases.MockForControllingIncomeData
	mockSapExportFailure      *mock_usecases.MockForLoggingSAPExportFailure
	mockIncomeFromTimesheet   *mock_usecases.MockForGettingIncomeFromTimesheet
}

// ExpectMirrorIncomeToTimesheet expects the income to be copied into income_from_timesheet as a
// brand new record, i.e. the timesheet consumer has not written anything for the period yet.
func (m *MockIncomeRepository) ExpectMirrorIncomeToTimesheet() {
	m.mockIncomeFromTimesheet.EXPECT().
		GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, ErrIncomeFromTimesheetNotFoundForPeriod)
	m.mockIncomeFromTimesheet.EXPECT().Add(gomock.Any()).Return(nil)
}

func (m *MockIncomeRepository) ExpectGetAllIncomeOfPreviousMonthByRole(incomes []*models.Income) {
	previousMonth := time.Now().AddDate(0, -1, 0)
	m.ExpectGetAllIncomeOfCurrentMonthByRole(incomes, previousMonth)
}

func (m *MockIncomeRepository) ExpectGetStudentLoans() {
	m.mockRead.EXPECT().GetStudentLoans().Return(models.StudentLoanList{List: []models.StudentLoan{}})
}

func (m *MockIncomeRepository) ExpectAddExport() {
	m.mockWrite.EXPECT().AddExport(gomock.Any()).Return(nil)
}

func (m *MockIncomeRepository) ExpectGetAllIncomeOfCurrentMonthByRole(incomes []*models.Income, now time.Time) {
	startDate, endDate := models.GetStartDateAndEndDate(now)
	m.mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		gomock.Any(), startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) ExpectGetAllIncomeByRoleStartDateAndEndDate(incomes []*models.Income, role string, startDate, endDate time.Time) {
	m.mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		role, startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) ExpectGetUserByID(id string) {
	m.mockGettingUserByID.EXPECT().GetByID(id).Return(&models.User{ID: bsonutil.MustObjectIDFromHex(id)}, nil)
}

func (m *MockIncomeRepository) ExpectGetCurrentUserIncomeNotFound(id string) {
	year, month := time.Now().Year(), time.Now().Month()
	m.mockControllingUserIncome.EXPECT().GetIncomeUserByYearMonth(id, year, month).Return(nil, errors.New("not found"))
}

func (m *MockIncomeRepository) ExpectAddIncomeSuccess() {
	m.mockControllingUserIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
}

func (m *MockIncomeRepository) ExpectGetIncomeByID(incID, uID string, income *models.Income) {
	m.mockUpdatingUserIncome.EXPECT().GetIncomeByID(incID, uID).Return(income, nil)
}

func (m *MockIncomeRepository) ExpectUpdateIncomeSuccess() {
	m.mockUpdatingUserIncome.EXPECT().UpdateIncome(gomock.Any()).Return(nil)
}

func deepClone(income *models.Income) *models.Income {
	b, _ := json.Marshal(income)
	var i models.Income
	json.Unmarshal(b, &i)
	return &i
}
