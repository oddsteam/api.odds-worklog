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
	ExportIncomeByStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) (string, error)
}
