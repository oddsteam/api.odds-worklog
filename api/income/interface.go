package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	AddIncome(u *models.Income) error
	GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error)
	GetIncomeByID(incID, uID string) (*models.Income, error)
	GetIncomeByUserID(uID string, fromYear int, fromMonth time.Month) (*models.Income, error)
	GetIncomeByStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) (*models.Income, error)
	UpdateIncome(income *models.Income) error
	AddExport(ep *models.Export) error
	GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error)
	UpdateExportStatus(id string) error
	GetStudentLoans() models.StudentLoanList
	SaveStudentLoans(loans models.StudentLoanList) int
	GetAllIncomeByStartDateAndEndDate(userIds []string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
	//GetAllUserIdByRole(role string) ([]string,error)

}

type Usecase interface {
	AddIncome(req *models.IncomeReq, uid string) (*models.Income, error)
	UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error)
	GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error)
	GetIncomeByUserIdAndCurrentMonth(userID string) (*models.Income, error)
	ExportIncome(role string, beforeMonth string) (string, error)
	ExportPdf(id string) (string, error)
	ExportIncomeNotExport(role string) (string, error)
	GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error)
	GetAllInComeByStartDateAndEndDate(userIds []string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
	ExportIncomeByStartDateAndEndDate(role string, incomes []*models.Income) (string, error)
	GetByRole(role string) ([]*models.User, error)
}
