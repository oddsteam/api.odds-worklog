package login

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Usecase interface {
	ValidateAndExtractToken(idToken string) (models.Identity, error)
	CreateUser(email string) (*models.User, error)
	CreateUserAndValidateEmail(email string) (*models.User, error)
	IsValidConsumerClientID(cid string) bool
}
