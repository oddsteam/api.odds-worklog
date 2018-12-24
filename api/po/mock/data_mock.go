package mock_po

import (
	"encoding/json"

	models "gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	Po = models.Po{
		ID:         bson.ObjectIdHex("5c1f855d59fc7d06988c6e01"),
		CustomerId: "5bbcf2f90fd2df527bc3c001",
		Name:       "KTB",
	}
	Po2 = models.Po{
		ID:         bson.ObjectIdHex("5c1f856459fc7d06988c6e02"),
		CustomerId: "5bbcf2f90fd2df527bc3c002",
		Name:       "AIS",
	}
	Poes = []*models.Po{&Po, &Po2}

	poByte, _   = json.Marshal(Po)
	PoJson      = string(poByte)
	poesByte, _ = json.Marshal(Poes)
	PoesJson    = string(poesByte)
)
