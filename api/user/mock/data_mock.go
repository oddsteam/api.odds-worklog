package mock_user

import (
	"encoding/json"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
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
		Vat:               "N",
		SlackAccount:      "test@abc.com",
		DailyIncome:       "5000",
		StatusTavi:        true,
		Address:           "every Where",
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
		DailyIncome:       "2000",
		StatusTavi:        true,
		Address:           "every Where",
	}

	StatusTavi = models.StatusTavi{
		ID:   bson.ObjectIdHex("5bbcf2f90fd2df527bc39539"),
		User: &User,
	}
	StatusTavi2 = models.StatusTavi{
		ID:   bson.ObjectIdHex("5bbcf2f90fd2df527bc39535"),
		User: &User2,
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
		DailyIncome:       "2000",
		StatusTavi:        true,
		Address:           "every Where",
	}

	adminByte, _ = json.Marshal(
		models.UserClaims{
			ID:         Admin.ID.Hex(),
			Role:       Admin.Role,
			StatusTavi: Admin.StatusTavi,
		},
	)
	AdminJson = string(adminByte)

	Token = models.Token{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY",
	}

	userByte, _ = json.Marshal(User)
	UserJson    = string(userByte)

	Users        = []*models.User{&User, &User2}
	ListUser     = []*models.StatusTavi{&StatusTavi}
	usersByte, _ = json.Marshal(Users)
	UsersJson    = string(usersByte)

	claimsUser = &models.JwtCustomClaims{
		&models.UserClaims{
			ID:         User.ID.Hex(),
			Role:       User.Role,
			StatusTavi: User.StatusTavi,
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	TokenUser = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsUser)

	claimsAdmin = &models.JwtCustomClaims{
		&models.UserClaims{
			ID:         Admin.ID.Hex(),
			Role:       Admin.Role,
			StatusTavi: Admin.StatusTavi,
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	TokenAdmin = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAdmin)
)
