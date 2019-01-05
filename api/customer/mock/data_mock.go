package mock_customer

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	Customer = models.Customer{
		ID:      bson.ObjectIdHex("5bbcf2f90fd2df527bc3c001"),
		Name:    "SEC",
		Address: "1234/123",
	}

	Customer2 = models.Customer{
		ID:      bson.ObjectIdHex("5bbcf2f90fd2df527bc3c002"),
		Name:    "SEC",
		Address: "1234/123",
	}

	cByte, _     = json.Marshal(Customer)
	CustomerJson = string(cByte)

	Customers     = []*models.Customer{&Customer, &Customer2}
	msByte, _     = json.Marshal(Customers)
	CustomersJson = string(msByte)
)
