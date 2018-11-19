package login

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	oauth2 "google.golang.org/api/oauth2/v2"
)

type Usecase interface {
	ManageLogin(idToken string) (*models.Token, error)
	GetTokenInfo(idToken string) (*oauth2.Tokeninfo, error)
	CreateUser(email string) (*models.User, error)
}
