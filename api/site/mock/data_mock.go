package mock_site

import (
	models "gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

var (
	MockSite = models.Site{
		ID:   bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539"),
		Name: "ktb",
	}
	MockSite2 = models.Site{
		ID:   bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39530"),
		Name: "ais",
	}
	MockSites = []*models.Site{&MockSite, &MockSite2}

	SiteJson = `{"name": "ktb"}`
)
