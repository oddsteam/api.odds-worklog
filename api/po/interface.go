package po

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	Create(po *models.Po) (*models.Po, error)
	Update(po *models.Po) (*models.Po, error)
}

type Usecase interface {
	Create(po *models.Po) (*models.Po, error)
	Update(po *models.Po) (*models.Po, error)
}
