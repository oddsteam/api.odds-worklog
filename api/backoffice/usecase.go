package backoffice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	repo     Repository
	siteRepo site.Repository
}

func NewUsecase(r Repository, sr site.Repository) Usecase {
	return &usecase{r, sr}
}

func (u *usecase) Get() ([]*models.UserIncome, error) {
	users, err := u.repo.Get()
	if err != nil {
		return nil, err
	}

	sites, err := u.siteRepo.GetSiteGroup()
	if err != nil {
		return nil, err
	}

	for _, s := range sites {
		for i, us := range users {
			if s.ID.Hex() == us.SiteID {
				users[i].Site = s
				users[i].SiteID = ""
				break
			}
		}
	}
	return users, nil
}