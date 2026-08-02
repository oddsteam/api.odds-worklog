package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo  Repository
	users ForGettingUsersBySiteID
}

func NewUsecase(r Repository, users ForGettingUsersBySiteID) Usecase {
	return &usecase{r, users}
}

func (u *usecase) CreateSiteGroup(m *models.Site) (*models.Site, error) {
	s, _ := u.repo.GetSiteGroupByName(m.Name)
	if s != nil {
		return nil, utils.ErrConflict
	}
	return u.repo.CreateSiteGroup(m)
}

func (u *usecase) UpdateSiteGroup(m *models.Site) (*models.Site, error) {
	existing, _ := u.repo.GetSiteGroupByName(m.Name)
	if existing != nil && existing.ID != m.ID {
		return nil, utils.ErrConflict
	}
	return u.repo.UpdateSiteGroup(m)
}

func (u *usecase) GetSiteGroup() ([]*models.Site, error) {
	return u.repo.GetSiteGroup()
}

func (u *usecase) GetSiteGroupByID(id string) (*models.Site, error) {
	return u.repo.GetSiteGroupByID(id)
}

func (u *usecase) DeleteSiteGroup(id string) error {
	users, err := u.users.GetBySiteID(id)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return utils.ErrSiteInUse
	}
	return u.repo.DeleteSiteGroup(id)
}
