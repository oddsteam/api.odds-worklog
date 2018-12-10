package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type Repository interface {
	CreateSiteGroup(sites *models.Site) (*models.Site, error)
	UpdateSiteGroup(sites *models.Site) (*models.Site, error)
	GetSiteGroup() ([]*models.Site, error)
	GetSiteGroupByID(id string) (*models.Site, error)
	DeleteSiteGroup(id string) error
}

type Usecase interface {
	CreateSiteGroup(m *models.Site) (*models.Site, error)
	UpdateSiteGroup(m *models.Site) (*models.Site, error)
	GetSiteGroup() ([]*models.Site, error)
	GetSiteGroupByID(id string) (*models.Site, error)
	DeleteSiteGroup(id string) error
}
