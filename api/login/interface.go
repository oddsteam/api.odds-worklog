package login

import (
	oauth2 "google.golang.org/api/oauth2/v2"
)

type Usecase interface {
	GetTokenInfo(idToken string) (*oauth2.Tokeninfo, error)
}
