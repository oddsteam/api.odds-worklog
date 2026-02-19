package mock_login

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var (
	MockToken = models.Token{
		Token:      "1234",
		FirstLogin: "Y",
	}

	LoginJson = `{"token": "5bbcf2f90fd2df527bc39539"}`
	Login     = models.Login{Token: "5bbcf2f90fd2df527bc39539"}
)
