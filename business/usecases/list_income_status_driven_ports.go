package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForGettingUsersByRole interface {
	GetByRole(role string) ([]*models.User, error)
}
