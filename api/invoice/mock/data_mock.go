package mock_invoice

import (
	"encoding/json"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"github.com/globalsign/mgo/bson"
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

	Invoice2 = models.Invoice{
		ID:        bson.ObjectIdHex("5bbcf2f90fd2df527bc30001"),
		PoID:      "5bbcf2f90fd2df527bcpo00",
		InvoiceNo: "2018_002",
		Amount:    "20000",
	}
	iByte2, _    = json.Marshal(Invoice2)
	InvoiceJson2 = string(iByte2)

	Invoices     = []*models.Invoice{&Invoice, &Invoice2}
	iBytes, _    = json.Marshal(Invoices)
	InvoicesJson = string(iBytes)

	InvoiceNoReqJson = `{"poID":"1234"}`
)
