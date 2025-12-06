package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/entity"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	AddIncome(u *models.Income) error
	AddIncomeOnSpecificTime(u *models.Income, t time.Time) error
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
	GetAllIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
	//GetAllUserIdByRole(role string) ([]string,error)

}

type Usecase interface {
	AddIncome(req *entity.IncomeReq, uid string) (*models.Income, error)
	UpdateIncome(id string, req *entity.IncomeReq, uid string) (*models.Income, error)
	GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error)
	GetIncomeByUserIdAndCurrentMonth(userID string) (*models.Income, error)
	ExportIncome(role string, beforeMonth string) (string, error)
	ExportIncomeByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error)
	ExportPdf(id string) (string, error)
	GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error)
	GetAllInComeByStartDateAndEndDate(userIds []string, startDate time.Time, endDate time.Time) ([]*models.Income, error)
	GetByRole(role string) ([]*models.User, error)
	ExportIncomeSAPByStartDateAndEndDate(role string, startDate, endDate time.Time, dateEff time.Time) (string, error)
}
