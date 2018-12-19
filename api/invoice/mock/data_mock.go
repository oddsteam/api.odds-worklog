package mock_invoice

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	Invoice = models.Invoice{
		ID:        bson.ObjectIdHex("5bbcf2f90fd2df527bc30000"),
		PoID:      "5bbcf2f90fd2df527bcpo00",
		InvoiceNo: "2018_001",
		Amount:    "10000",
	}
	iByte, _    = json.Marshal(Invoice)
	InvoiceJson = string(iByte)
)
