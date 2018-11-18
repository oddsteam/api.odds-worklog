package mock_login

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	oauth2 "google.golang.org/api/oauth2/v2"
)

var (
	MockTokenInfo = oauth2.Tokeninfo{
		Audience: "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com",
		Email:    "abc@mail.com",
	}

	LoginJson = `{"token": "5bbcf2f90fd2df527bc39539"}`
	Login     = models.Login{Token: "5bbcf2f90fd2df527bc39539"}
)
