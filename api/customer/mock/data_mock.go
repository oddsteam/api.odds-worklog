package mock_user

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	MockCustomer = models.Customer{
		ID:      bson.ObjectIdHex("5bbcf2f90fd2df527bc39539"),
		Name:    "SEC",
		Address: "1234/123",
	}

	MockToken = models.Token{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY",
	}

	MockCustomerById = models.Customer{
		ID:      "1234567890",
		Name:    "DTAC",
		Address: "1234/123",
	}

	MockCustomerById2 = models.Customer{
		ID:      "1234567891",
		Name:    "SEC",
		Address: "1234/123",
	}

	MockAdmin = models.User{
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

	userByte, _ = json.Marshal(MockCustomer)
	UserJson    = string(userByte)

	MockCustomers   = []*models.Customer{&MockCustomerById, &MockCustomerById2}
	UserListByte, _ = json.Marshal(MockCustomers)
	UserListJson    = string(UserListByte)
)
