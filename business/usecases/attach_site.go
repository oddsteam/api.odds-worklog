package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

// attachSite resolves user.SiteID into user.Site so payroll can snapshot SiteName.
// Missing site or empty SiteID leaves Site unset; write paths must not fail for that.
func attachSite(user *models.User, siteRepo ForGettingSiteByID) {
	if user == nil || siteRepo == nil || user.SiteID == "" {
		return
	}
	site, err := siteRepo.GetSiteGroupByID(user.SiteID)
	if err != nil || site == nil {
		return
	}
	user.Site = site
}
