package consumer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	Create(u *models.Consumer) (*models.Consumer, error)
	GetByClientID(cid string) (*models.Consumer, error)
}

type Usecase interface {
	GetByClientID(cid string) (*models.Consumer, error)
}
