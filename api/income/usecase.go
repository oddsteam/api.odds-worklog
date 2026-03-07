package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

func NewUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
}
