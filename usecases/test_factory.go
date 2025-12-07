package usecases

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/file"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/usecases/mock"
)

func CreateExportIncomeUsecaseWithMock(t *testing.T) (ForUsingExportIncome, *gomock.Controller, *MockIncomeRepository) {
	ctrl := gomock.NewController(t)
	mockRepoIncome := mockIncomeRepository(ctrl)

	usecase := NewExportIncomeUsecase(mockRepoIncome.mockRead, mockRepoIncome.mockWrite, userMock.NewMockRepository(ctrl), file.NewCSVWriter(), file.NewSAPWriter())
	return usecase, ctrl, mockRepoIncome
}

func mockIncomeRepository(ctrl *gomock.Controller) *MockIncomeRepository {
	mockRepoIncome := MockIncomeRepository{
		mock_usecases.NewMockForGettingIncomeData(ctrl),
		mock_usecases.NewMockForControllingIncomeData(ctrl)}
	mockRepoIncome.mockWrite.EXPECT().AddExport(gomock.Any()).Return(nil).AnyTimes()
	mockRepoIncome.mockRead.EXPECT().GetStudentLoans().AnyTimes()
	return &mockRepoIncome
}

type MockIncomeRepository struct {
	mockRead  *mock_usecases.MockForGettingIncomeData
	mockWrite *mock_usecases.MockForControllingIncomeData
}

func (m *MockIncomeRepository) ExpectGetAllIncomeOfPreviousMonthByRole(incomes []*models.Income) {
	previousMonth := time.Now().AddDate(0, -1, 0)
	m.ExpectGetAllIncomeOfCurrentMonthByRole(incomes, previousMonth)
}

func (m *MockIncomeRepository) ExpectGetAllIncomeOfCurrentMonthByRole(incomes []*models.Income, now time.Time) {
	startDate, endDate := utils.GetStartDateAndEndDate(now)
	m.mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		gomock.Any(), startDate, endDate).Return(incomes, nil)
}

func (m *MockIncomeRepository) GetAllIncomeByRoleStartDateAndEndDate(incomes []*models.Income, role string, startDate, endDate time.Time) {
	m.mockRead.EXPECT().GetAllIncomeByRoleStartDateAndEndDate(
		role, startDate, endDate).Return(incomes, nil)
}

func deepClone(income *models.Income) *models.Income {
	b, _ := json.Marshal(income)
	var i models.Income
	json.Unmarshal(b, &i)
	return &i

}
