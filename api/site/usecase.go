package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) CreateSiteGroup(m *models.Site) (*models.Site, error) {
	s, _ := u.repo.GetSiteGroupByName(m.Name)
	if s != nil {
		return nil, utils.ErrConflict
	}
	return u.repo.CreateSiteGroup(m)
}

func (u *usecase) UpdateSiteGroup(m *models.Site) (*models.Site, error) {
	return u.repo.UpdateSiteGroup(m)
}

func (u *usecase) GetSiteGroup() ([]*models.Site, error) {
	return u.repo.GetSiteGroup()
}

func (u *usecase) GetSiteGroupByID(id string) (*models.Site, error) {
	return u.repo.GetSiteGroupByID(id)
}

func (u *usecase) DeleteSiteGroup(id string) error {
	return u.repo.DeleteSiteGroup(id)
}
