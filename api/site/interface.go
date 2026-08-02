package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

type Repository interface {
	CreateSiteGroup(site *models.Site) (*models.Site, error)
	UpdateSiteGroup(site *models.Site) (*models.Site, error)
	GetSiteGroup() ([]*models.Site, error)
	GetSiteGroupByID(id string) (*models.Site, error)
	GetSiteGroupByName(name string) (*models.Site, error)
	DeleteSiteGroup(id string) error
}

// ForGettingUsersBySiteID looks up users assigned to a site (avoids importing the user package).
type ForGettingUsersBySiteID interface {
	GetBySiteID(id string) ([]*models.User, error)
}

type Usecase interface {
	CreateSiteGroup(m *models.Site) (*models.Site, error)
	UpdateSiteGroup(m *models.Site) (*models.Site, error)
	GetSiteGroup() ([]*models.Site, error)
	GetSiteGroupByID(id string) (*models.Site, error)
	DeleteSiteGroup(id string) error
}
