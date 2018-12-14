package mock_site

import (
	models "gitlab.odds.team/worklog/api.odds-worklog/models"
	"gopkg.in/mgo.v2/bson"
)

var (
	MockSite = models.Site{
		ID:   bson.ObjectIdHex("5bbcf2f90fd2df527bc39539"),
		Name: "ktb",
	}
	MockSite2 = models.Site{
		ID:   bson.ObjectIdHex("5bbcf2f90fd2df527bc39530"),
		Name: "ais",
	}
	MockSites = []*models.Site{&MockSite, &MockSite2}

	SiteJson = `{"name": "ktb"}`
)
