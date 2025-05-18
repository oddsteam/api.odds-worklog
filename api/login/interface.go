package login

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	oauth2 "google.golang.org/api/oauth2/v2"
)

type Usecase interface {
	GetTokenInfo(idToken string) (*oauth2.Tokeninfo, error)
	CreateUserAndValidateEmail(email string) (*models.User, error)
	IsValidConsumerClientID(cid string) bool
}
