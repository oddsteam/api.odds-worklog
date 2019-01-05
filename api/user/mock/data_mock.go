package mock_user

import (
	"encoding/json"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	User = models.User{
		ID:                bson.ObjectIdHex("5bbcf2f90fd2df527bc39539"),
		Role:              "corporate",
		FirstName:         "Tester",
		LastName:          "Super",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		Vat:               "Y",
		SlackAccount:      "test@abc.com",
	}

	User2 = models.User{
		ID:                bson.ObjectIdHex("5bbcf2f90fd2df527bc39530"),
		Role:              "corporate",
		FirstName:         "Tester",
		LastName:          "Super",
		Email:             "test@abc.com",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		Vat:               "Y",
		SlackAccount:      "test@abc.com",
	}

	Admin = models.User{
		ID:                bson.ObjectIdHex("5bbcf2f90fd2df527bc39535"),
		Role:              "admin",
		FirstName:         "Tester",
		LastName:          "Super",
		Email:             "jin@odds.team",
		BankAccountName:   "ทดสอบชอบลงทุน",
		BankAccountNumber: "123123123123",
		ThaiCitizenID:     "1234567890123",
		Vat:               "Y",
		SlackAccount:      "test@abc.com",
	}

	adminByte, _ = json.Marshal(Admin)
	AdminJson    = string(adminByte)

	Token = models.Token{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY",
	}

	userByte, _ = json.Marshal(User)
	UserJson    = string(userByte)

	Users        = []*models.User{&User, &User2}
	usersByte, _ = json.Marshal(Users)
	UsersJson    = string(usersByte)

	claimsUser = &models.JwtCustomClaims{
		&User,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	TokenUser = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsUser)

	claimsAdmin = &models.JwtCustomClaims{
		&Admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	TokenAdmin = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAdmin)
)
