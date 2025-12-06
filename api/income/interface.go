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
	GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error)
	UpdateExportStatus(id string) error
	SaveStudentLoans(loans models.StudentLoanList) int
}

type Usecase interface {
	AddIncome(req *entity.IncomeReq, uid string) (*models.Income, error)
	UpdateIncome(id string, req *entity.IncomeReq, uid string) (*models.Income, error)
	GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error)
	GetIncomeByUserIdAndCurrentMonth(userID string) (*models.Income, error)
	ExportPdf(id string) (string, error)
	GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error)
	GetByRole(role string) ([]*models.User, error)
}
