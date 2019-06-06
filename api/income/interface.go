package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	AddIncome(u *models.Income) error
	GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error)
	GetIncomeByID(incID, uID string) (*models.Income, error)
	GetIncomeByUserID(uID string) (*models.Income, error)
	UpdateIncome(income *models.Income) error
	AddExport(ep *models.Export) error
	UpdateExportStatus(id string) error
}

type Usecase interface {
	AddIncome(req *models.IncomeReq, user *models.User) (*models.Income, error)
	UpdateIncome(id string, req *models.IncomeReq, user *models.User) (*models.Income, error)
	GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error)
	GetIncomeByUserIdAndCurrentMonth(userID string) (*models.Income, error)
	ExportIncome(role string, beforeMonth string) (string, error)
	ExportPdf(id string) (string, error)
	ExportIncomeNotExport(role string) (string, error)
}
