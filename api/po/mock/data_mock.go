package mock_po

import (
	"encoding/json"

	models "gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	MockPo = models.Po{
		ID:         bson.ObjectIdHex("5c1f855d59fc7d06988c6e01"),
		CustomerId: "01",
		Name:       "ktb",
	}
	MockPo2 = models.Po{
		ID:         bson.ObjectIdHex("5c1f856459fc7d06988c6e02"),
		CustomerId: "02",
		Name:       "ais",
	}
	MockPoes = []*models.Po{&MockPo, &MockPo2}

	PoByte, _  = json.Marshal(MockPo)
	MockPoJson = string(PoByte)
)
