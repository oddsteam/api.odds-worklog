package usecases

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/golang/mock/gomock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/usecases/mock"
)

func CreateExportIncomeUsecaseWithMock(t *testing.T) (ForUsingExportIncome, *gomock.Controller, *MockIncomeRepository) {
	ctrl := gomock.NewController(t)
	mockRepoIncome := mockIncomeRepository(ctrl)

	usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, file.NewCSVWriter(), file.NewSAPWriter())
	return usecase, ctrl, mockRepoIncome
}

func CreateAddIncomeUsecaseWithMock(mockRepoIncome *MockIncomeRepository) ForUsingAddIncome {
	usecase := NewAddIncomeUsecase(mockRepoIncome.mockControllingUserIncome, mockRepoIncome.mockGettingUserByID)
	return usecase
}

func mockIncomeRepository(ctrl *gomock.Controller) *MockIncomeRepository {
	mockRepoIncome := MockIncomeRepository{
		mock_usecases.NewMockForGettingUserByID(ctrl),
		mock_usecases.NewMockForControllingUserIncome(ctrl),
		mock_usecases.NewMockForGettingIncomeData(ctrl),
		mock_usecases.NewMockForControllingIncomeData(ctrl)}
	return &mockRepoIncome
}

type MockIncomeRepository struct {
	mockGettingUserByID       *mock_usecases.MockForGettingUserByID
	mockControllingUserIncome *mock_usecases.MockForControllingUserIncome
	mockRead                  *mock_usecases.MockForGettingIncomeData
	mockWrite                 *mock_usecases.MockForControllingIncomeData
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
	startDate, endDate := utils.GetStartDateAndEndDate(now)
	m.mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		gomock.Any(), startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) ExpectGetAllIncomeByRoleStartDateAndEndDate(incomes []*models.Income, role string, startDate, endDate time.Time) {
	m.mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		role, startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) ExpectGetUserByID(id string) {
	m.mockGettingUserByID.EXPECT().GetByID(id).Return(&models.User{ID: bson.ObjectIdHex(id)}, nil)
}

func (m *MockIncomeRepository) ExpectGetCurrentUserIncomeNotFound(id string) {
	year, month := time.Now().Year(), time.Now().Month()
	m.mockControllingUserIncome.EXPECT().GetIncomeUserByYearMonth(id, year, month).Return(nil, errors.New("not found"))
}

func (m *MockIncomeRepository) ExpectAddIncomeSuccess() {
	m.mockControllingUserIncome.EXPECT().AddIncome(gomock.Any()).Return(nil)
}

func deepClone(income *models.Income) *models.Income {
	b, _ := json.Marshal(income)
	var i models.Income
	json.Unmarshal(b, &i)
	return &i
}
